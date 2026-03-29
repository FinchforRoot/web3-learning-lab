![img_9.png](img_9.png)
代码
```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RoleManage{
    enum Role { Owner,Admin,User,None}
    mapping ( address => Role) public userRole;
    address public owner;

    constructor(){
        owner = msg.sender;
        userRole[msg.sender] = Role.Owner;
    }

    modifier onlyOwner(){
        require(userRole[msg.sender] == Role.Owner,"Only owner can call");
        _;
    }

    modifier onlyAdmin(){
        require(userRole[msg.sender] == Role.Admin || userRole[msg.sender] == Role.Owner,"Only Admin and Owner can call");
        _;
    }

    function addAdmin(address user) public onlyOwner {
        // Owner添加Admin
        userRole[user] = Role.Admin;
    }
    
    function addUser(address user) public onlyAdmin {
        // Admin添加User
        userRole[user] = Role.User;
    }
    
    function getRole(address user) public view returns (Role) {
        // 查询角色
        return userRole[user];
    }
    
}
```
问题1：
两者主要区别在于内部调用、gas成本以及参数类型。
public可以内部调用，但是external不可以；public更费gas；public的参数只能是memory，但是external可以是calldata

当只需要外部调用的时候||入参包含大数组或者长字符串时||需要优化gas时。

问题2：
view主要用在查看状态或者获取数据上
pure主要用在纯计算任务上，尤其是用作工具函数

问题3：
按照从左到右依次执行
例如
function test()public modifier1 modifier2 modifier3 {
        code
}
那么先执行modifier1再是modifier2，再是3，最后是code的内容