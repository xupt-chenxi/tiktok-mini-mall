// 网关入口
// Author: chenxi 2025.01
package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"tiktok-mini-mall/internal/app/gateway/handler"
	"tiktok-mini-mall/pkg/middleware"
	"tiktok-mini-mall/pkg/utils"
	"time"
)

func main() {
	utils.InitViper("configs/config.yaml")
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middleware.AuthMiddleware())
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
	cartGroup := r.Group("/cart")
	{
		cartGroup.POST("/add-item", handler.AddItemHandler)
		cartGroup.POST("/get-cart", handler.GetCartHandler)
		cartGroup.POST("/empty-cart", handler.EmptyCartHandler)
	}

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
