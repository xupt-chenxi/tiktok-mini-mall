package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"tiktok-mini-mall/internal/app/user/handler"
	"tiktok-mini-mall/internal/app/user/repository"
)

func main() {
	// 初始化 viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	dsn := viper.GetString("user.database.dsn")
	repository.InitDatabase(dsn)

	r := gin.Default()
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", handler.RegisterHandler)
		userGroup.POST("/login", handler.LoginHandler)
	}
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
