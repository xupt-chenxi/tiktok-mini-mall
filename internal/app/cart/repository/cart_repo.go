package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"sync"
	"tiktok-mini-mall/api/pb/cart"
	"tiktok-mini-mall/internal/app/cart/model"
)

var (
	db   *gorm.DB
	once sync.Once
)

// InitDatabase 初始化数据库连接
func InitDatabase(dsn string) {
	once.Do(func() { // 确保只初始化一次
		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("连接cart数据库失败: %v\n", err)
		}
		log.Println("cart数据库连接成功")
	})
	err := db.AutoMigrate(&model.Cart{})
	if err != nil {
		log.Printf("cart数据库自动迁移失败: %v\n", err)
	}
}

func AddItem(item *cart.CartItem, userId int64) error {
	cartItem := &model.Cart{
		UserId:    userId,
		ProductId: item.ProductId,
		Quantity:  item.Quantity,
	}
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "product_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"quantity": gorm.Expr("quantity + ?", cartItem.Quantity),
		}),
	}).Create(cartItem).Error
}

func GetCart(userId int64) ([]*model.Cart, error) {
	var cartItems []*model.Cart
	result := db.Where("user_id = ?", userId).Find(&cartItems)
	if result.Error != nil {
		return nil, result.Error
	}
	return cartItems, nil
}

func EmptyCart(userId int64) error {
	result := db.Where("user_id = ?", userId).Delete(&model.Cart{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
