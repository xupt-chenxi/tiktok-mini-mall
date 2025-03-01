package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"tiktok-mini-mall/api/pb/user"
	"tiktok-mini-mall/pkg/utils"
)

var (
	userClient user.UserServiceClient
	userConn   *grpc.ClientConn
	userMu     sync.Mutex
)

func GetUserClient() (user.UserServiceClient, error) {
	if userClient == nil || !isUserConnHealthy() {
		userMu.Lock()
		defer userMu.Unlock()
		if userConn != nil {
			_ = userConn.Close()
		}
		if userClient == nil {
			ip, port := utils.Config.User.IP, utils.Config.User.Port
			var err error
			userConn, err = grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return nil, err
			}
			log.Printf("已建立用户服务的gRPC连接")
			userClient = user.NewUserServiceClient(userConn)
		}
	}

	return userClient, nil
}

func isUserConnHealthy() bool {
	state := userConn.GetState()
	return state == connectivity.Ready || state == connectivity.Idle
}
