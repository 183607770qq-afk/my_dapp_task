// main.go
package main

import (
    "task/database"
    "task/listener"
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/joho/godotenv"
)

func main() {
    // 加载环境变量
    godotenv.Load()
    
    // 初始化数据库
    database.InitDB()

    // 从环境变量获取配置
    rpcURL := os.Getenv("ETH_RPC_URL")
    if rpcURL == "" {
        rpcURL = "ws://localhost:8545" // 默认本地节点
    }
    log.Printf("rpcUrl",rpcURL)
    contractAddress := os.Getenv("CONTRACT_ADDRESS")
    if contractAddress == "" {
        log.Fatal("CONTRACT_ADDRESS environment variable is required")
    }

    // 创建事件监听器
    eventListener, err := listener.NewEventListener(rpcURL, contractAddress)
    if err != nil {
        log.Fatalf("Failed to create event listener: %v", err)
    }

    // 设置优雅关闭
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 启动监听器
    go func() {
        if err := eventListener.Start(ctx); err != nil {
            log.Fatalf("Event listener failed: %v", err)
        }
    }()

    log.Println("Contract event listener started successfully")
    log.Printf("Listening to contract: %s", contractAddress)

    // 等待中断信号
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh

    log.Println("Shutting down event listener...")
    cancel()
    
    // 给清理工作一些时间
    time.Sleep(2 * time.Second)
    log.Println("Event listener stopped")
}