package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"tiktok-mini-mall/api/pb/user"
	"tiktok-mini-mall/pkg/utils"
)

var userClient user.UserServiceClient

// var userConn *grpc.ClientConn
var userMu sync.Mutex

func GetUserClient() (user.UserServiceClient, error) {
	if userClient == nil {
		userMu.Lock()
		defer userMu.Unlock()
		if userClient == nil {
			ip, port := utils.Config.User.IP, utils.Config.User.Port
			userConn, err := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return nil, err
			}
			log.Printf("已建立用户服务的gRPC连接")
			userClient = user.NewUserServiceClient(userConn)
		}
	}

	return userClient, nil
}

//func CloseUserConn() error {
//	if userConn != nil {
//		err := userConn.Close()
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
