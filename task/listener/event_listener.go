// listener/event_listener.go
package listener

import (
	"context"
	"fmt"
	"log"
	"task/database"
	"task/models"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

type EventListener struct {
	client          *ethclient.Client
	db              *gorm.DB
	contractAddr    common.Address
	eventSignatures map[string]string // 事件签名映射
}

func NewEventListener(rpcURL, contractAddress string) (*EventListener, error) {
	// 连接以太坊节点，需要使用WebSocket协议[2](@ref)
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %v", err)
	}

	db := database.DB

	return &EventListener{
		client:       client,
		db:           db,
		contractAddr: common.HexToAddress(contractAddress),
		eventSignatures: map[string]string{
			// 常见ERC20事件签名[1](@ref)
			crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")).Hex(): "Transfer",
			crypto.Keccak256Hash([]byte("Approval(address,address,uint256)")).Hex(): "Approval",
		},
	}, nil
}

// 开始监听事件
func (el *EventListener) Start(ctx context.Context) error {
	log.Println("Starting event listener...")

	// 创建事件过滤查询[1,2](@ref)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{el.contractAddr},
	}

	logs := make(chan types.Log)

	// 订阅事件[2,6](@ref)
	sub, err := el.client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %v", err)
	}

	log.Println("Successfully subscribed to contract events")

	for {
		select {
		case err := <-sub.Err():
			log.Printf("Subscription error: %v", err)
			// 实现简单的重连逻辑
			el.handleReconnect(ctx)
		case vLog := <-logs:
			go el.processLog(vLog) // 并发处理事件
		case <-ctx.Done():
			log.Println("Stopping event listener...")
			return ctx.Err()
		}
	}
}

// 处理日志事件
func (el *EventListener) processLog(vLog types.Log) {
	eventName := el.getEventName(vLog.Topics[0].Hex())

	// 解析事件数据
	eventData, err := el.parseLogData(vLog)
	if err != nil {
		log.Printf("Error parsing log data: %v", err)
		return
	}

	// 创建事件记录
	event := models.ContractEvent{
		TransactionHash: vLog.TxHash.Hex(),
		BlockNumber:     vLog.BlockNumber,
		ContractAddress: vLog.Address.Hex(),
		EventName:       eventName,
		EventData:       eventData,
		LogIndex:        vLog.Index,
		CreatedAt:       time.Now(),
	}

	// 保存到数据库
	result := el.db.Create(&event)
	if result.Error != nil {
		log.Printf("Error saving event to database: %v", result.Error)
		return
	}

	log.Printf("Event saved: %s from block %d, tx: %s",
		eventName, vLog.BlockNumber, vLog.TxHash.Hex())
}
