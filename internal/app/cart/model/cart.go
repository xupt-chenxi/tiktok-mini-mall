package model

type Cart struct {
	Id        uint32 `gorm:"primaryKey"`
	UserId    int64  `gorm:"type:int;not null;index:idx_user_product,unique"`
	ProductId uint32 `gorm:"type:int;not null;index:idx_user_product,unique"`
	Quantity  int32  `gorm:"type:int;not null"`
}
