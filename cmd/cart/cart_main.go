package main

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	cart "tiktok-mini-mall/api/pb/cart_pb"

	"tiktok-mini-mall/internal/app/cart/repository"
	"tiktok-mini-mall/internal/app/cart/service"
	"tiktok-mini-mall/pkg/utils"
)

func main() {
	utils.InitViper("configs/config.yaml")
	dsn := viper.GetString("cart.database.dsn")
	repository.InitDatabase(dsn)

	port := viper.GetString("cart.port")
	// 注册购物车品服务
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("cart服务监听端口失败: %v", err)
	}
	grpcServer := grpc.NewServer()
	cart.RegisterCartServiceServer(grpcServer, &service.CartService{})
	log.Println("cart服务启动...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("cart服务启动失败: %v", err)
	}
}
