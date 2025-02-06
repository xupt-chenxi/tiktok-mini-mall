package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
	"tiktok-mini-mall/internal/app/shop/model"
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
			log.Printf("连接order数据库失败: %v\n", err)
		}
		log.Println("order数据库连接成功")
	})
	err := db.AutoMigrate(&model.Order{})
	if err != nil {
		log.Printf("order数据库自动迁移失败: %v\n", err)
	}
}

func AddOrder(order *model.Order) error {
	return db.Create(order).Error
}

func GetListOrder(userId int64) ([]*model.Order, error) {
	var orders []*model.Order
	result := db.Where("user_id = ?", userId).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

func MarkOrderPaid(userId int64, orderId string) error {
	return db.Model(&model.Order{}).Where("id = ? AND user_id = ?", orderId, userId).Update("state", 1).Error
}
