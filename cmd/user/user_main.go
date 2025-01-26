package main

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	userpb "tiktok-mini-mall/api/pb/user_pb"
	"tiktok-mini-mall/internal/app/user/repository"
	"tiktok-mini-mall/internal/app/user/service"
	"tiktok-mini-mall/pkg/utils"
)

func main() {
	utils.InitViper("configs/config.yaml")
	dsn := viper.GetString("user.database.dsn")
	repository.InitDatabase(dsn)

	port := viper.GetString("user.port")
	// 注册用户服务
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("用户服务监听端口失败: %v", err)
	}
	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &service.UserService{})
	log.Println("用户服务启动...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("用户服务启动失败: %v", err)
	}
}
