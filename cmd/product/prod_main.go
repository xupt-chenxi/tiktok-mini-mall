package main

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	prod "tiktok-mini-mall/api/pb/prod_pb"
	"tiktok-mini-mall/internal/app/product/repository"
	"tiktok-mini-mall/internal/app/product/service"
	"tiktok-mini-mall/pkg/utils"
)

func main() {
	utils.InitViper("configs/config.yaml")
	dsn := viper.GetString("product.database.dsn")
	repository.InitDatabase(dsn)

	port := viper.GetString("product.port")
	// 注册商品服务
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("商品服务监听端口失败: %v", err)
	}
	grpcServer := grpc.NewServer()
	prod.RegisterProductCatalogServiceServer(grpcServer, &service.ProductService{})
	log.Println("商品服务启动...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("商品服务启动失败: %v", err)
	}
}
