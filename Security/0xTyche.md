---
timezone: UTC+8
---

> 请在上边的 timezone 添加你的当地时区(UTC)，这会有助于你的打卡状态的自动化更新，如果没有添加，默认为北京时间 UTC+8 时区


# 你的名字

1. 自我介绍
0xTyche 对合约安全、defi十分感兴趣
2. 你认为你会完成这次共学小组吗？
希望会
3. 你感兴趣的小组
合约安全
4. 你的联系方式（Wechat or Telegram）
0xTyche

## Notes

<!-- Content_START -->

### 2025.11.18
Part I - 动手部署一个智能合约 writeup
1. 完成对报名  
2. 测试网络领水：https://www.alchemy.com/faucets/ethereum-sepolia
address：0x00000000bb09009cdcd358d6c5ce6f56611577f1  
![image](https://github.com/0xTyche/PKUBA-Colearn-25-Fall/blob/main/pictures/get-sepolia-eth.png)  
3. 登录remix 网站：https://remix.ethereum.org/

```Solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract HelloWeb3 {
    event Greeting(address indexed sender, uint256 timestamp);
    
    constructor() {}

    function hello() external {
        emit Greeting(msg.sender, block.timestamp);
    }
}
```
4. 合约部署成功 tx：https://web3.okx.com/zh-hans/explorer/sepolia/tx/0xb8dd671fa5bc78ba53150ac40fc3591c17d86e053de763572230521d7b25d026
  
合约地址：0xa7120cc8d48f4053c2eb0babb449d20f2ab9af49

```shell
[block:9655856 txIndex:46]from: 0x000...577f1to: HelloWeb3.(constructor)value: 0 weidata: 0x608...00033logs: 0hash: 0x074...9af58
view on Etherscan view on Blockscout
Verification process started...
Verifying with Sourcify...
Verifying with Routescan...
Etherscan verification skipped: API key not found in global Settings.
Sourcify verification successful.
https://repo.sourcify.dev/11155111/0xA7120Cc8D48F4053c2eb0BaBb449d20f2Ab9Af49/
Routescan verification successful.
https://testnet.routescan.io/address/0xA7120Cc8D48F4053c2eb0BaBb449d20f2Ab9Af49/contract/11155111/code
```


### 2025.07.11

笔记内容

### 2025.07.12

<!-- Content_END -->
