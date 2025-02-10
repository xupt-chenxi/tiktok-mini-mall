package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
	"tiktok-mini-mall/internal/app/user/model"
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
			log.Printf("user连接数据库失败: %v\n", err)
		} else {
			log.Println("user数据库连接成功")
		}
	})
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		log.Printf("数据库自动迁移失败: %v\n", err)
	}
}

func CreateUser(user *model.User) error {
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := db.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
