package main

import (
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
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

	namingClient, err := utils.NewNamingClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          utils.Config.Shop.IP,
		Port:        50054,
		ServiceName: "shop_service",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	if success == false {
		log.Fatalf("shop服务注册失败")
	}
	log.Println("shop服务注册成功")

	port := utils.Config.Shop.Port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("shop服务监听端口失败: %v", err)
	}
	grpcServer := grpc.NewServer()
	shop.RegisterShopServiceServer(grpcServer, &service.ShopService{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("shop服务启动失败: %v", err)
	}
}
