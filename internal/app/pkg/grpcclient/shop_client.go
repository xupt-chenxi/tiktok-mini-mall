package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"tiktok-mini-mall/api/pb/shop"
	"tiktok-mini-mall/pkg/utils"
)

var (
	shopClient shop.ShopServiceClient
	shopConn   *grpc.ClientConn
	shopMu     sync.Mutex
)

func GetShopClient() (shop.ShopServiceClient, error) {
	if shopClient == nil || !isShopConnHealthy() {
		shopMu.Lock()
		defer shopMu.Unlock()
		if shopConn != nil {
			_ = shopConn.Close()
		}
		if shopClient == nil {
			ip, port := utils.Config.Shop.IP, utils.Config.Shop.Port
			var err error
			shopConn, err = grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return nil, err
			}
			log.Printf("已建立Shop服务的gRPC连接")
			shopClient = shop.NewShopServiceClient(shopConn)
		}
	}

	return shopClient, nil
}

func isShopConnHealthy() bool {
	state := shopConn.GetState()
	return state == connectivity.Ready || state == connectivity.Idle
}
