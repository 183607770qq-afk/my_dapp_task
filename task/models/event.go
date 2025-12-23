package models

import (
	"time"
)

type ContractEvent struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TransactionHash string    `gorm:"size:66;not null;index" json:"transaction_hash"`
	BlockNumber     uint64    `gorm:"not null;index" json:"block_number"`
	ContractAddress string    `gorm:"size:42;not null;index" json:"contract_address"`
	EventName       string    `gorm:"size:100;not null;index" json:"event_name"`
	EventData       string    `gorm:"type:text" json:"event_data"` // JSON格式数据
	LogIndex        uint      `gorm:"not null" json:"log_index"`
	CreatedAt       time.Time `json:"created_at"`
}
