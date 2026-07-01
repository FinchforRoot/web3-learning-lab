### 练习1：创建基础事件

**任务**：创建一个简单的留言板合约，使用事件记录所有留言。

**要求**：
1. 定义MessagePosted事件，包含：用户地址、留言内容、时间戳
2. 实现postMessage函数，触发事件
3. 正确使用indexed参数

```
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Message {
    
    event MessagePosted(address indexed account,string msg,uint256 timestamp);

    function postMessage(string memory message)external {
        require(bytes(message).length > 0,"Message cannot be empty");
        require(bytes(message).length < 100,"Message too long");
        emit MessagePosted(msg.sender, message, block.timestamp);
    }
}
```

### 练习2：实现代币事件

**任务**：创建一个简单的ERC20代币合约，实现Transfer和Approval事件。

**要求**：
1. 实现Transfer事件（包含from、to、value）
2. 实现Approval事件（包含owner、spender、value）
3. 在transfer、approve、transferFrom函数中正确触发事件
4. 正确处理铸造（from=0）和销毁（to=0）的情况

```

```