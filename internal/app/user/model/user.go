package model

import "time"

type User struct {
	Id        int64  `gorm:"primaryKey;autoIncrement:false"`
	Email     string `gorm:"uniqueIndex;not null;size:255"`
	PassHash  string `gorm:"not null;size:64"` // 加密后的密码
	Nickname  string `gorm:"size:100"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//Orders []Order `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` // 关联订单
	//Cart   []Cart  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` // 关联购物车
}
