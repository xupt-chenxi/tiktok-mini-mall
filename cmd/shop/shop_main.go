package main

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	shop "tiktok-mini-mall/api/pb/shop_pb"
	"tiktok-mini-mall/internal/app/shop/repository"
	"tiktok-mini-mall/internal/app/shop/service"
	"tiktok-mini-mall/pkg/utils"
)

func main() {
	utils.InitViper("configs/config.yaml")
	dsn := viper.GetString("shop.database.dsn")
	repository.InitDatabase(dsn)

	port := viper.GetString("shop.port")
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
