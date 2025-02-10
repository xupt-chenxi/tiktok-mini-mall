package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"tiktok-mini-mall/api/pb/cart"
	"tiktok-mini-mall/internal/app/cart/repository"
	"tiktok-mini-mall/internal/app/cart/service"
	"tiktok-mini-mall/pkg/utils"
)

func main() {
	dsn := utils.Config.Cart.Database.DSN
	repository.InitDatabase(dsn)

	port := utils.Config.Cart.Port
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
