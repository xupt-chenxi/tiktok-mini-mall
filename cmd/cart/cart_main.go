package main

import (
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
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

	namingClient, err := utils.NewNamingClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          utils.Config.Cart.IP,
		Port:        50053,
		ServiceName: "cart_service",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	if success == false {
		log.Fatalf("cart服务注册失败")
	}
	log.Println("cart服务注册成功")

	port := utils.Config.Cart.Port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("cart服务监听端口失败: %v", err)
	}
	grpcServer := grpc.NewServer()
	cart.RegisterCartServiceServer(grpcServer, &service.CartService{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("cart服务启动失败: %v", err)
	}
}
