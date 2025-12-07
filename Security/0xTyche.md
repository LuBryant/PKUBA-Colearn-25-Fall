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

Part II - 智能合约编写

题目：通过编写智能合约与靶子合约交互，获取 Flag 并触发 ChallengeCompleted 事件。

https://sepolia.etherscan.io/address/0x4a6C0c0dc8bD8276b65956c9978ef941C3550A1B#code

根据查看合约代码，可以知道hint函数所给的暗示是"keccak PKUBlockchain"，后续还是尽量按照题目要求逐步完成。

1. 首先调用目标合约中hint函数，看题目给了什么暗示
```shell
root@racknerd-9da1d08:~/home/PKUBA-1/PKUBA-Colearn-25-Fall/writeup/part2# cast call $TARGET_CONTRACT "hint()(string)" --rpc-url $SEPOLIA_RPC_URL
"keccak PKUBlockchain"
```
2. 因此我们需要使用keccak单向加密函数对PKUBlockchain字符串进行加密，加密的结果可能就是答案。
3. 解答项目结构
/src/Solver.sol  
```solidity
// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

interface ITargetContract {
    function hint() external pure returns (string memory);
    function query(bytes32 _hash) external returns (string memory);
}

contract Solver {
    ITargetContract public target;
    
    constructor(address _targetAddress) {
        target = ITargetContract(_targetAddress);
    }
    
    // 获取提示
    function getHint() external view returns (string memory) {
        // 调用目标合约中的hint函数
        return target.hint();
    }
    
    // 解答
    function getCorrectHash() public pure returns (bytes32) {
        return keccak256(abi.encodePacked("PKUBlockchain"));
    }
    
    // 提交答案
    function solve() external returns (string memory) {
        bytes32 answer = getCorrectHash();
        return target.query(answer);
    }
}
```

script/Deploy.s.sol  
```Solidity
// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Script.sol";
import "../src/Solver.sol";

contract DeployScript is Script {
    function run() external {
        address targetContract = vm.envAddress("TARGET_CONTRACT");
        
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        
        vm.startBroadcast(deployerPrivateKey);
        
        // 部署 Solver 合约
        Solver solver = new Solver(targetContract);
        
        console.log("Solver deployed at:", address(solver));
        console.log("Target contract:", targetContract);
        
        // 自动调用 solve 方法
        string memory flag = solver.solve();
        console.log("flag result", flag);
        
        vm.stopBroadcast();
    }
}
```
run
```shell
forge script script/Deploy.s.sol:DeployAndSolve \
    --rpc-url $SEPOLIA_RPC_URL \
    --broadcast \
    -vvvv

Script ran successfully.

== Logs ==
  Solver deployed at: 0x6F00A229cf51DB7Eec4B6996F2eBcFE365C0Ae98
  Target contract: 0x4a6C0c0dc8bD8276b65956c9978ef941C3550A1B
  flag result FLAG{PKU_Blockchain_Colearn_Week1_Success}
```
3. 看到上方合约执行的结果可以知道，我们成功的完成了题目的要求。
0x6F00A229cf51DB7Eec4B6996F2eBcFE365C0Ae98 这个是我部署合约的地址，
按要求通过合约调用  
https://sepolia.etherscan.io/tx/0x91e72d0a469e800d7f44f2a02b40518128a5a59eea8124e85496997113082604

### 2025.07.11



### 2025.12.07

转账哈希：0x6ea04b5764b8db4cc59f7f3f872a45df6fcc0b9d1b8345c4725786de052d0051

本周目标：

学会用 Geth 的 Go 客户端从 RPC 节点读取链上信息
理解区块链底层数据结构（block、transaction、receipt）
这些是做 DApp、链上分析、合约调试等的基础。

2. 节点：运行客户端的软件实例
从功能和数据完整性角度分析，节点可以分为
- 全节点：保存当前完整状态和必要的历史数据，可以独立验证新区块和交易；
- 轻节点：只保存少量数据和区块头，需要向其他节点请求详细信息；
- 归档节点：不仅保存当前状态，还保存所有历史状态，方便做历史查询和分析，但资源消耗很大。

有矿工/验证者节点把交易打包进入区块，交易才算真正的上链。

3. rpc 给外部程序调用的接口
RPC（Remote Procedure Call，远程过程调用）是节点向外暴露的一组标准接口，用来让其他程序查询或提交数据，常见是 HTTP RPC 或 WebSocket RPC。

Geth（Go Ethereum）是用 Go 语言实现的以太坊客户端，也是目前使用最广泛的实现之一。
【看来Go学习也要提上日程了】

公共的Sepolia rpc：https://ethereum-sepolia-rpc.publicnode.com

Go 客户端库
Geth 提供了一套Go客户端库，方便Go代码和以太坊交互，常用的是ethclient包。
```Go
    import (
        "context"
        "github.com/ethereum/go-ethereum/ethclient"
    )
    func main(){
        // 连接公共 rpc
        client, err := ethclient.Dial("https://ethereum-sepolia-rpc.publicnode.com")

        // 也可以使用个人节点
        // client, err := ethclient.Dial("https://your-server-ip:8545") 和 anvil fork 出一个节点类似

        if err != nil {
            panic(err)
        }
        defer client.Close()
    }
```

<!-- Content_END -->
