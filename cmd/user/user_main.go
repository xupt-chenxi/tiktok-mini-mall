package main

import (
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
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

	namingClient, err := utils.NewNamingClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          utils.Config.User.IP,
		Port:        50051,
		ServiceName: "user_service",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	if success == false {
		log.Fatalf("用户服务注册失败")
	}
	log.Println("用户服务注册成功")

	port := utils.Config.User.Port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("用户服务监听端口失败: %v", err)
	}
	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, &service.UserService{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("用户服务启动失败: %v", err)
	}
}
