package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ERC-20 Transfer 事件签名
var transferEventSignature = common.HexToHash(
	"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
)

type TxTracker struct {
	TxHash         string          // 查询的交易hash
	Status         string          // 交易的状态：pending、success、fail
	BlockNumber    uint64          // 交易所在的区块
	From           common.Address  // 发送者地址
	To             *common.Address // 接收者地址【部署合约时 to 是 nil，所以用指针】
	Value          *big.Int        // ETH数量，单位 wei
	GasUsed        uint64          // 实际消耗的 Gas 量
	GasPrice       *big.Int        // 单位 Gas 价格（wei），EIP-1559 后是有效价格
	GasFee         *big.Int        // Gas 总费用
	FromBalance    *big.Int        // 交易发出方当前 ETH 余额
	ToBalance      *big.Int        // 交易接收方当前 ETH 余额（wei）
	TokenTransfers []TokenTransfer // 这笔交易里触发的所有 ERC-20 Transfer 事件
}

type TokenTransfer struct {
	TokenAddress string
	Symbol       string
	From         string
	To           string
	Amount       *big.Int
}

func main() {
	rpcURL := os.Getenv("ETH_RPC_URL")
	if rpcURL == "" {
		log.Fatal("ETH_RPC_URL is not set")
	}
	txHashHex := flag.String("tx", "", "transaction hash (for query mode)")
	flag.Parse()
	// 1. 连接节点
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum node: %v", err)
	}

	// 2. 输入交易哈希
	txHash := common.HexToHash(*txHashHex)

	// 3. 追踪交易
	tracker, err := TrackTransaction(client, txHash)
	if err != nil {
		panic(err)
	}

	// 4. 打印报告
	PrintReport(tracker)
}

func TrackTransaction(client *ethclient.Client, txHash common.Hash) (*TxTracker, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// 先获取交易信息
	tx, ispending, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		// 如果找不到，可能是无效哈希
		return nil, fmt.Errorf("交易未找到: %v", err)
	}
	// 根据交易hash获取发送者的地址
	from, err := getFrom(client, tx)
	if err != nil {
		return nil, fmt.Errorf("获取发送方失败: %v", err)
	}

	// 先构造一个交易追踪结构体的返回值
	tracker := &TxTracker{
		TxHash: txHash.Hex(),
		From:   from,
	}

	if ispending {
		tracker.Status = "pending"
		return tracker, nil
	}

	// 根据交易hash获取回执信息
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		tracker.Status = "fail"
		return tracker, fmt.Errorf("获取回执失败且交易不在pending池: %v", err)
	}

	if receipt == nil {
		tracker.Status = "fail"
		return tracker, fmt.Errorf("获取的回执信息为nil,稍后再重试")
	}

	tracker.Status = "success"
	tracker.BlockNumber = receipt.BlockNumber.Uint64()
	tracker.To = tx.To()                         // *common.Address，可能为 nil（部署合约时）
	tracker.Value = tx.Value()                   // *big.Int，单位 wei
	tracker.GasPrice = receipt.EffectiveGasPrice // 实际支付的单位 gas 价格
	tracker.GasUsed = receipt.GasUsed            // uint64
	// GasFee = GasUsed × GasPrice
	tracker.GasFee = new(big.Int).Mul(
		big.NewInt(int64(receipt.GasUsed)),
		receipt.EffectiveGasPrice,
	)
	// 发送方余额
	fromBalance, err := client.BalanceAt(ctx, tracker.From, nil)
	if err != nil {
		return nil, fmt.Errorf("获取发送方余额失败: %v", err)
	}
	tracker.FromBalance = fromBalance

	// 接收方余额（如果 to 存在）
	if tx.To() != nil {
		toBalance, err := client.BalanceAt(ctx, *tx.To(), nil)
		if err != nil {
			return nil, fmt.Errorf("获取接收方余额失败: %v", err)
		}
		tracker.ToBalance = toBalance
	}
	// 代币转账
	tracker.TokenTransfers = parseTokenTransfers(receipt)
	return tracker, nil
}

func parseTokenTransfers(receipt *types.Receipt) []TokenTransfer {
	var transfers []TokenTransfer

	for _, vLog := range receipt.Logs {
		// Transfer 事件的 topics[0] 是这个签名
		if len(vLog.Topics) == 3 && vLog.Topics[0] == transferEventSignature {
			transfer := TokenTransfer{
				TokenAddress: vLog.Address.Hex(),
				From:         common.HexToAddress(vLog.Topics[1].Hex()).Hex(),
				To:           common.HexToAddress(vLog.Topics[2].Hex()).Hex(),
				Amount:       new(big.Int).SetBytes(vLog.Data),
			}
			transfers = append(transfers, transfer)
		}
	}
	return transfers
}

func getFrom(client *ethclient.Client, tx *types.Transaction) (common.Address, error) {
	ctx := context.Background()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return common.Address{}, fmt.Errorf("获取链ID失败: %v", err)
	}
	from, err := types.Sender(types.LatestSignerForChainID(chainID), tx)
	if err != nil {
		return common.Address{}, fmt.Errorf("恢复发送方失败: %v", err)
	}
	return from, nil
}

func PrintReport(tracker *TxTracker) {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("🔍 交易追踪报告")
	fmt.Println(strings.Repeat("=", 60))

	// 基本信息
	fmt.Printf("交易哈希: %s\n", tracker.TxHash)
	fmt.Printf("状态:      %s\n", formatStatus(tracker.Status))
	if tracker.BlockNumber > 0 {
		fmt.Printf("区块号:    %d\n", tracker.BlockNumber)
	}
	fmt.Printf("发送方:    %s\n", tracker.From.Hex())

	if tracker.To != nil {
		fmt.Printf("接收方:    %s\n", tracker.To.Hex())
	} else {
		fmt.Println("接收方:    合约部署 (to 为空)")
	}

	// ETH 转账金额
	if tracker.Value != nil && tracker.Value.Cmp(big.NewInt(0)) > 0 {
		fmt.Printf("转账金额:  %s ETH\n", weiToEther(tracker.Value))
	} else {
		fmt.Println("转账金额:  0 ETH")
	}

	fmt.Println(strings.Repeat("-", 60))

	// Gas 信息
	if tracker.GasUsed > 0 {
		fmt.Printf("Gas 消耗:  %d\n", tracker.GasUsed)
	}
	if tracker.GasPrice != nil {
		fmt.Printf("Gas 价格:  %s gwei\n", weiToGwei(tracker.GasPrice))
	}
	if tracker.GasFee != nil {
		fmt.Printf("Gas 总费用: %s ETH\n", weiToEther(tracker.GasFee))
	}

	fmt.Println(strings.Repeat("-", 60))

	// 余额信息
	if tracker.FromBalance != nil {
		fmt.Printf("发送方余额: %s ETH\n", weiToEther(tracker.FromBalance))
	}
	if tracker.ToBalance != nil {
		fmt.Printf("接收方余额: %s ETH\n", weiToEther(tracker.ToBalance))
	}

	// 代币转账
	if len(tracker.TokenTransfers) > 0 {
		fmt.Println(strings.Repeat("-", 60))
		fmt.Printf("🪙 代币转账 (%d 笔):\n", len(tracker.TokenTransfers))
		for i, transfer := range tracker.TokenTransfers {
			fmt.Printf("\n  [%d] 代币合约: %s\n", i+1, transfer.TokenAddress)
			if transfer.Symbol != "" {
				fmt.Printf("      代币符号: %s\n", transfer.Symbol)
			}
			fmt.Printf("      发送方:   %s\n", transfer.From)
			fmt.Printf("      接收方:   %s\n", transfer.To)
			if transfer.Amount != nil {
				fmt.Printf("      数量:     %s\n", transfer.Amount.String())
			}
		}
	}

	fmt.Println(strings.Repeat("=", 60))
}

// 状态美化
func formatStatus(status string) string {
	switch status {
	case "success":
		return "✅ 成功"
	case "fail":
		return "❌ 失败"
	case "pending":
		return "⏳ 等待中"
	default:
		return status
	}
}

// wei 转 ETH
func weiToEther(wei *big.Int) string {
	ether := new(big.Float).Quo(
		new(big.Float).SetInt(wei),
		new(big.Float).SetInt(big.NewInt(1e18)),
	)
	return fmt.Sprintf("%.6f", ether)
}

// wei 转 gwei
func weiToGwei(wei *big.Int) string {
	gwei := new(big.Float).Quo(
		new(big.Float).SetInt(wei),
		new(big.Float).SetInt(big.NewInt(1e9)),
	)
	return fmt.Sprintf("%.2f", gwei)
}
