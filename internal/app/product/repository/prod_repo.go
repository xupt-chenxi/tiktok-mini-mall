package repository

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
	"tiktok-mini-mall/internal/app/product/model"
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
			log.Printf("连接数据库失败: %v\n", err)
		}
		log.Println("数据库连接成功")
	})
	err := db.AutoMigrate(&model.Product{})
	if err != nil {
		log.Printf("数据库自动迁移失败: %v\n", err)
	}
}

func GetProductList(page int, pageSize int, categoryName string) ([]*model.Product, error) {
	var productList []*model.Product
	offset := (page - 1) * pageSize
	var result *gorm.DB
	if categoryName != "" {
		result = db.Where("JSON_CONTAINS(categories, ?)", fmt.Sprintf(`["%s"]`, categoryName)).Offset(offset).Limit(pageSize).Find(&productList)
	} else {
		result = db.Offset(offset).Limit(pageSize).Find(&productList)
	}
	return productList, result.Error
}

func GetProductById(id int) (*model.Product, error) {
	product := &model.Product{}
	result := db.First(&product, id)
	return product, result.Error
}

func SearchProducts(query string) ([]*model.Product, error) {
	var products []*model.Product
	// TODO 先基于模糊匹配进行搜索, 后续引入 Elasticsearch
	result := db.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").Find(&products)
	return products, result.Error
}
