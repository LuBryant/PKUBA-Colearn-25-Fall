package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// ------------------------------------------------
// ⚠️ 关键配置：本地 WebSocket 节点地址
// ------------------------------------------------

const (
	// 本地 Geth 节点的 WebSocket 地址
	// 默认端口：8546 (WebSocket), 8545 (HTTP RPC)
	// 确保你的本地 Geth 节点已启动并启用了 WebSocket
	NodeWSS = "ws://127.0.0.1:8546"

	// 连接超时时间
	CONNECTION_TIMEOUT = 30 * time.Second
)

func main() {
	log.Println("开始连接到本地 WebSocket 节点")

	// 1. 建立底层的 RPC 连接 (WebSocket)
	// 注意：必须用 rpc.DialContext 建立基础连接，以便复用
	ctx, cancel := context.WithTimeout(context.Background(), CONNECTION_TIMEOUT)
	defer cancel()

	rpcClient, err := rpc.DialContext(ctx, NodeWSS)
	if err != nil {
		log.Fatalf("❌ 无法连接到本地 WebSocket 节点: %v\n"+
			"   可能的原因：\n"+
			"   1. 本地 Geth 节点未启动\n"+
			"   2. WebSocket 未启用（需要在启动 Geth 时添加 --ws 参数）\n"+
			"   3. 端口配置错误（默认 WebSocket 端口为 8546）\n"+
			"   提示：启动 Geth 节点示例：geth --ws --ws.addr 0.0.0.0 --ws.port 8546", err)
	}
	defer rpcClient.Close()
	fmt.Println("✅ 成功建立 RPC WebSocket 连接")

	// 3. 初始化两种不同的 Client
	// EthClient: 用于通用查询和区块头订阅
	ethClient := ethclient.NewClient(rpcClient)
	// GethClient: 用于 Geth 特有的订阅 (如 Pending Transactions)
	gethClient := gethclient.New(rpcClient)

	// 4. 创建数据通道
	newHeadChan := make(chan *types.Header) // 接收新区块头
	pendingTxChan := make(chan common.Hash)  // 接收 Pending 交易 Hash

	// 5. 开启订阅
	// A. 订阅新区块 (SubscribeNewHead)
	headSub, err := ethClient.SubscribeNewHead(context.Background(), newHeadChan)
	if err != nil {
		log.Fatalf("❌ 订阅新区块失败: %v", err)
	}
	fmt.Println("🎧 开始监听新区块 (NewHeads)...")

	// B. 订阅待处理交易 (SubscribePendingTransactions)
	// 注意：本地 Geth 节点完全支持此功能
	txSub, err := gethClient.SubscribePendingTransactions(context.Background(), pendingTxChan)
	if err != nil {
		log.Printf("⚠️  警告: 订阅 Pending 交易失败: %v\n"+
			"   可能的原因：\n"+
			"   1. Geth 节点版本过旧，不支持此功能\n"+
			"   2. 节点配置问题\n"+
			"   建议：检查 Geth 版本和配置", err)
		// 继续运行，只监听区块
		txSub = nil
	} else {
		fmt.Println("🎧 开始监听交易池 (Pending Transactions)...")
	}

	// 6. 优雅退出信号捕获
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// 7. 主循环：处理接收到的数据
	fmt.Println("\n📡 监控已启动，按 Ctrl+C 退出...\n")
	for {
		select {
		// 处理新区块
		case header := <-newHeadChan:
			fmt.Printf("\n📦 [New Block] Height: %d | Hash: %s | Time: %d\n",
				header.Number, header.Hash().Hex(), header.Time)

			// 实际应用场景：在这里触发你的业务逻辑，例如检查 Uniswap 价格

		// 处理 Pending 交易
		case txHash := <-pendingTxChan:
			// 为了演示不刷屏，我们只打印 Hash，实际中你会在这里并发去 fetch 交易详情
			fmt.Printf("🌊 [Pending Tx] %s\n", txHash.Hex())

			// 模拟 MEV 逻辑：
			// go analyzeTransaction(ethClient, txHash)

		// 处理订阅错误 (如网络断开)
		case err := <-headSub.Err():
			log.Fatalf("❌ 区块订阅异常中断: %v", err)
		case err := <-txSub.Err():
			if txSub != nil {
				log.Fatalf("❌ 交易订阅异常中断: %v", err)
			}

		// 用户退出
		case <-sigChan:
			fmt.Println("\n🛑 停止监控，正在断开连接...")
			headSub.Unsubscribe()
			if txSub != nil {
				txSub.Unsubscribe()
			}
			return
		}
	}
}

// 模拟分析函数 (伪代码)
func analyzeTransaction(client *ethclient.Client, hash common.Hash) {
	// tx, isPending, err := client.TransactionByHash(context.Background(), hash)
	// 1. 解码 Input Data 看是不是在调用 Uniswap Router
	// 2. 模拟执行看利润
	// 3. 发送 Bundle
}

