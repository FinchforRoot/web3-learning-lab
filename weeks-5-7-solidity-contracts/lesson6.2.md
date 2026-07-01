### 练习1：创建数学库
**任务**：

编写一个Solidity库合约，实现基本的数学运算。

**要求**：

1. 实现平方根函数（使用Newton-Raphson方法）
2. 实现最大公约数（GCD）
3. 实现幂运算
4. 所有函数都是pure函数

```
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library SafeMath {
    function sqrt(uint x) internal pure returns (uint) {
        if (x == 0) return 0;
        uint z = (x + 1) / 2;
        uint y = x;
        while (z < y) {
            y = z;
            z = (x / z + z) / 2;
        }
        return y;
    }

    function gcd(uint256 a, uint256 b) internal pure returns (uint256) {
        while (b != 0) {
            uint256 temp = b;
            b = a % b;
            a = temp;
        }
        return a;
    }

    function power(
        uint256 base,
        uint256 exponent
    ) internal pure returns (uint256) {
        if (exponent == 0) return 1;
        uint256 result = 1;
        uint256 currentBase = base;
        while (exponent > 0) {
            if (exponent % 2 == 1) {
                result *= currentBase;
            }
            exponent /= 2;
            currentBase *= currentBase;
        }
        return result;
    }
}

```

### 练习2：使用using for改写合约

```
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library CounterLib{
    function increment(uint256 value) internal pure returns (uint256) {
        return value + 1;
    }
    
    function decrement(uint256 value) internal pure returns (uint256) {
        require(value > 0, "Cannot decrement zero");
        return value - 1;
    }
    
    function reset(uint256 /* value */) internal pure returns (uint256) {
        return 0;
    }
}

contract Counter {
    using CounterLib for uint256;
    
    uint256 public count;
    
    function increment() public {
        count = count.increment();
    }
    
    function decrement() public {
        count = count.decrement();
    }
    
    function reset() public {
        count = count.reset();
    }
}
```

### 练习3：地址白名单库

EnumerableSet的核心思想是组合两种数据结构的优势：
- array提供遍历能力
- mapping提供O(1)查找能力
- 同步维护两个结构
- 删除时使用交换技巧

```
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library AddressSet{
    struct Set{
        address[] values;
        mapping(address => uint256) indexs;
    }

    function contains(Set storage set, address value) internal view returns (bool){
        return set.indexs[value] != 0;
    }

    function length(Set storage set) internal view returns(uint256){
        return set.values.length;
    }

    // 返回参数1的参数2的地址
    function at(Set storage set,uint256 index)internal view returns(address){
        require(index < set.values.length,"index out of bounds");
        return set.values[index];
    }

    function add(Set storage set, address value) internal returns (bool){
        if (contains(set, value)){
            return false;
        }
        set.values.push(value);
        set.indexs[value] = set.values.length;
        return true;
    }

    function remove(Set storage set,address value)internal returns(bool){
        //1.先查找位置
        uint index = set.indexs[value];
        if (index == 0) {
            return false;
        }
        // 记录下要删除的index
        uint256 toDeleteIndex = index -1;
        // 记录最后一个数组元素的index
        uint256 lastIndex = set.values.length -1;
        // 如果不是最后一个元素，那么就把数组最后一个元素的值和数组要删除元素的值做个对调，同时更新原来最后一个元素的索引下标位置
        if (toDeleteIndex != lastIndex) {
            address lastAddress = set.values[lastIndex];
            set.values[toDeleteIndex] = lastAddress;
            set.indexs[lastAddress] = index;
        }
        set.values.pop();
        // 最后再把mapping的key为原来value进行delete
        delete set.indexs[value];
        return true;
    }
}

contract WhiteList{
    using AddressSet for AddressSet.Set;

    AddressSet.Set private whiteList;

    function addToWhiteList(address account) public {
        require(whiteList.add(account),"Already in whitelist");
    }

    function removeFromWhiteList(address account) public {
        require(whiteList.remove(account),"Not in whitelist");
    }

    function isWhiteListed(address addr)public view  returns (bool){
        return whiteList.contains(addr);
    }

    function getWhiteListSize()public view returns (uint256){
        return whiteList.length();
    }

}
```