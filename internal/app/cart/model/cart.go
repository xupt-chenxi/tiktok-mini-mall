package model

type Cart struct {
	Id        uint32 `gorm:"primaryKey"`
	UserId    int64  `gorm:"type:int;not null"`
	ProductId uint32 `gorm:"uniqueIndex;type:int;not null"`
	Quantity  int32  `gorm:"type:int;not null"`
}
