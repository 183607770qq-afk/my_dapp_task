// listener/event_parser.go
package listener

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (el *EventListener) getEventName(topicHash string) string {
    if name, exists := el.eventSignatures[topicHash]; exists {
        return name
    }
    return "Unknown"
}

func (el *EventListener) parseLogData(vLog types.Log) (string, error) {
    eventData := make(map[string]interface{})
    
    // 解析Transfer事件（示例）[1](@ref)
    if len(vLog.Topics) >= 3 && el.getEventName(vLog.Topics[0].Hex()) == "Transfer" {
        eventData["from"] = common.BytesToAddress(vLog.Topics[1].Bytes()).Hex()
        eventData["to"] = common.BytesToAddress(vLog.Topics[2].Bytes()).Hex()
        
        // 解析金额数据（在Data字段中）
        if len(vLog.Data) >= 32 {
            // 简化处理：将数据转换为大整数
            value := common.BytesToHash(vLog.Data).Big()
            eventData["value"] = value.String()
        }
    }
    
    // 可以在这里添加其他事件类型的解析逻辑
    
    jsonData, err := json.Marshal(eventData)
    if err != nil {
        return "", err
    }
    
    return string(jsonData), nil
}

// 重连处理
func (el *EventListener) handleReconnect(ctx context.Context) {
    log.Println("Attempting to reconnect...")
    // 实现简单的重连逻辑
    time.Sleep(5 * time.Second)
}