// 网关入口
// Author: chenxi 2025.01
package main

import (
	"github.com/gin-gonic/gin"
	"tiktok-mini-mall/internal/app/gateway/handler"
	"tiktok-mini-mall/pkg/middleware"
	"tiktok-mini-mall/pkg/utils"
)

func main() {
	utils.InitViper("configs/config.yaml")
	r := gin.Default()
	r.Use(middleware.TraceIDMiddleware())
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", handler.RegisterHandler)
		userGroup.POST("/login", handler.LoginHandler)
	}
	productGroup := r.Group("/product")
	{
		productGroup.GET("/list-products", handler.ListProductsHandler)
		productGroup.GET("/:id", handler.GetProductHandler)
		productGroup.GET("/search", handler.SearchProductsHandler)
	}

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
