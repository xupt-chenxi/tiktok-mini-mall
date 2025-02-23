// 网关入口
// Author: chenxi 2025.01
package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
	"tiktok-mini-mall/internal/app/gateway/handler"
	"tiktok-mini-mall/pkg/middleware"
	"tiktok-mini-mall/pkg/utils"
	"time"
)

func main() {
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
	shopGroup := r.Group("/shop")
	{
		shopGroup.POST("/place-order", handler.PlaceOrderHandler)
		shopGroup.POST("/list-order", handler.ListOrderHandler)
		shopGroup.POST("/update-order", handler.UpdateOrderState)
	}

	namingClient, err := utils.NewNamingClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          utils.Config.Gateway.IP,
		Port:        8080,
		ServiceName: "gateway_service",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	if success == false {
		log.Fatalf("网关服务注册失败")
	}
	log.Println("网关服务注册成功")
	err = r.Run(utils.Config.Gateway.Port)
	if err != nil {
		log.Fatalf("网关服务启动失败: %v", err)
	}
}
