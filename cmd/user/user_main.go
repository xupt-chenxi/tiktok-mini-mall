package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"tiktok-mini-mall/api/pb/user"
	"tiktok-mini-mall/internal/app/user/repository"
	"tiktok-mini-mall/internal/app/user/service"
	"tiktok-mini-mall/pkg/utils"
)

func main() {
	dsn := utils.Config.User.Database.DSN
	repository.InitDatabase(dsn)

	port := utils.Config.User.Port
	// 注册用户服务
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("用户服务监听端口失败: %v", err)
	}
	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, &service.UserService{})
	log.Println("用户服务启动...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("用户服务启动失败: %v", err)
	}
}
