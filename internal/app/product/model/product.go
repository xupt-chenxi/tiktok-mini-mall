package model

type Product struct {
	ID          uint32  `gorm:"primaryKey"`
	Name        string  `gorm:"type:varchar(255);not null"`
	Description string  `gorm:"type:text"`
	Picture     string  `gorm:"type:varchar(255)"`
	Price       float32 `gorm:"type:float;not null"`
	Categories  string  `gorm:"type:json"` // 使用 JSON 类型存储数组
	Stock       uint32  `gorm:"type:int;not null"`
}
