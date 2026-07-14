package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	rpcURL := os.Getenv("ETH_RPC_URL")
	if rpcURL == "" {
		log.Fatal("ETH_RPC_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("failed to get chain id: %v", err)
	}
	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║     Ethereum 区块监控工具 v1.0      ║")
	fmt.Println("╠══════════════════════════════════════╣")
	fmt.Printf("║ RPC URL  : %-24s ║\n", rpcURL)
	fmt.Printf("║ Chain ID : %-24s ║\n", chainID.String())
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Println("🔍 开始监控新区块... (按 Ctrl+C 退出)\n")

	var lastBlockNumber uint64

	// 获取当前最新的区块数
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("获取最新区块失败: %v", err)
	}
	lastBlockNumber = header.Number.Uint64()

	// 获取完整的区块信息
	block, err := client.BlockByNumber(ctx, nil)
	if err == nil {
		printDetailedBlock(block, true)
	}
	// 3秒轮训一次
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			latestHeader, err := client.HeaderByNumber(ctx, nil)
			if err != nil {
				log.Printf("获取区块失败: %v", err)
				continue
			}
			currentBlock := latestHeader.Number.Uint64()
			if currentBlock > lastBlockNumber {
				for blockNum := lastBlockNumber + 1; blockNum < currentBlock; blockNum++ {
					block, err := client.BlockByNumber(ctx, big.NewInt(int64(blockNum)))
					if err != nil {
						log.Printf("获取区块 %d 失败: %v", blockNum, err)
						continue
					}
					printDetailedBlock(block, false)
				}
				lastBlockNumber = currentBlock
			} else {
				fmt.Printf("[%s] ⏳ 等待新区块... 当前高度: %d\r",
					time.Now().Format("15:04:05"), currentBlock)
			}
		case sig := <-sigChan:
			fmt.Printf("\n\n👋 收到信号 %v，监控结束，最终区块高度: %d\n", sig, lastBlockNumber)
			return
		}

	}
}

func printDetailedBlock(block *types.Block, isInitial bool) {
	timestamp := time.Unix(int64(block.Time()), 0)
	prefix := "🆕"
	if isInitial {
		prefix = "📍"
	}

	fmt.Printf("\n%s 区块 #%d\n", prefix, block.Number().Uint64())
	fmt.Printf("   ├─ 哈希     : %s\n", block.Hash().Hex())
	fmt.Printf("   ├─ 时间     : %s\n", timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("   ├─ 交易数   : %d\n", len(block.Transactions()))
	fmt.Printf("   ├─ Gas限制  : %d\n", block.GasLimit())
	fmt.Printf("   ├─ Gas使用  : %d\n", block.GasUsed())
	fmt.Printf("   └─ 难度     : %d\n", block.Difficulty().Uint64())
}
