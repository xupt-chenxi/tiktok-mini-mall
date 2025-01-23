package repository

import (
	"errors"
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
			log.Fatalf("连接数据库失败: %v", err)
		}
		log.Println("数据库连接成功")
	})
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("数据库自动迁移失败: %v", err)
	}
}

func CreateUser(user *model.User) error {
	result := db.Create(user)
	if result.Error != nil {
		log.Println("用户创建失败")
		return errors.New("用户创建失败")
	}
	return nil
}

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := db.First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("email 为: %s 的用户不存在\n", email)
			return nil, errors.New("用户不存在")
		}
	}
	return &user, nil
}
