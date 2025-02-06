package model

import "time"

type Order struct {
	Id         int64   `gorm:"primaryKey;autoIncrement:false"`
	UserId     int64   `gorm:"not null"`
	Name       string  `gorm:"type:varchar(100);not null"`
	Email      string  `gorm:"type:varchar(255);not null"`
	Address    string  `gorm:"type:varchar(255);not null"`
	Price      float32 `gorm:"not null"`
	OrderItems string  `gorm:"type:json;not null"` // 使用 JSON 类型存储数组
	State      uint8   `gorm:"not null"`           // 订单状态, 0: 未支付 1: 已支付 2: 已失效
	CreatedAt  time.Time
}
