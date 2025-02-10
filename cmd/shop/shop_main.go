package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"tiktok-mini-mall/api/pb/shop"
	"tiktok-mini-mall/internal/app/shop/repository"
	"tiktok-mini-mall/internal/app/shop/service"
	"tiktok-mini-mall/pkg/utils"
)

func main() {
	dsn := utils.Config.Shop.Database.DSN
	repository.InitDatabase(dsn)

	port := utils.Config.Shop.Port
	// 注册购物车品服务
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("shop服务监听端口失败: %v", err)
	}
	grpcServer := grpc.NewServer()
	shop.RegisterShopServiceServer(grpcServer, &service.ShopService{})
	log.Println("shop服务启动...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("shop服务启动失败: %v", err)
	}
}
