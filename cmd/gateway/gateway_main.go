// 网关入口
// Author: chenxi 2025.01
package main

import (
	"github.com/gin-gonic/gin"
	"tiktok-mini-mall/internal/app/gateway/handler"
	"tiktok-mini-mall/pkg"
)

func main() {
	pkg.InitViper("configs/config.yaml")
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
