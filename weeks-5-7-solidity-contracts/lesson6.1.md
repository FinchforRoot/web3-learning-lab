### 练习1：实现完整的权限管理系统

**需求**：

创建一个模块化的权限管理系统：

1. Ownable合约：单一所有者管理
2. Pausable合约：暂停功能
3. MyContract：组合两个功能

```
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;


contract Ownable {
    address public owner;

    event OwnershipTransferred(address indexed previousOwner,address indexed newOwner);

    constructor(){
        owner = msg.sender;
        emit OwnershipTransferred(address(0), msg.sender);
    }

    modifier onlyOwner(){
        require(msg.sender == owner,"only owner");
        _;
    }

    function transferOwnership(address newOwner) public onlyOwner{
        require(newOwner != address(0),"invalid address");
        owner = newOwner;
        emit OwnershipTransferred(owner,newOwner);
    }
 
}

contract Pausable {
    bool public paused;

    event Paused(address account);

    event UnPaused(address account);

    modifier whenPaused(){
        require(paused,"constract is paused!");
        _;
    }

    modifier whenNotPaused(){
        require(!paused,"constract is not paused!");
        _;
    }

    function _pause() internal whenNotPaused{
        paused = true;
        emit Paused(msg.sender);
    }

    function _unpaused() internal whenPaused{
        paused = false;
        emit UnPaused(msg.sender);
    }
 
}

contract MyContract is Ownable,Pausable{
    uint public value;

    function setVlaue(uint _value)public onlyOwner whenNotPaused{
        value = _value;
    }

    function pause()public onlyOwner{
        _pause();
    }

    function unpause() public onlyOwner{
        _unpaused();
    }

}
```

### 练习2：实现动物抽象合约

**需求**：

1. 创建Animal抽象合约，定义makeSound抽象函数
2. 创建Dog和Cat子合约实现makeSound
3. 添加共同的eat函数


```
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

abstract contract Animal {
    string public species;

    constructor(string memory _species){
        species = _species;
    }

    function makeSound() public pure virtual returns (string memory);

    function eat() public pure returns (string memory) {
        return "Eating...";
    }
    
    function sleep() public pure returns (string memory) {
        return "Sleeping...";
    }
}

contract Cat is Animal{

    constructor() Animal("Cat"){}

    function makeSound() public pure override  returns (string memory){
        return "miao miao";
    }

}

contract Dog is Animal{

    constructor() Animal("Dog"){}

    function makeSound() public pure override  returns (string memory){
        return "wang  wang";
    }

}
```

### 练习3：使用OpenZeppelin创建代币

**需求**：

使用OpenZeppelin库创建一个完整的代币合约：

1. 继承ERC20
2. 继承Ownable
3. 添加mint功能

```
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract MyToken is Ownable,ERC20{
    constructor(uint initialSupply) Ownable(msg.sender) ERC20("my token","MTK"){
        _mint(msg.sender,initialSupply);
    }

    function mint(address to, uint amount)public onlyOwner{
        _mint(to, amount);
    }
}
```