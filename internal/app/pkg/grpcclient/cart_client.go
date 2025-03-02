package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"tiktok-mini-mall/api/pb/cart"
	"tiktok-mini-mall/pkg/utils"
)

var (
	cartClient cart.CartServiceClient
	cartConn   *grpc.ClientConn
	cartMu     sync.Mutex
)

func GetCartClient() (cart.CartServiceClient, error) {
	if cartClient == nil || !isCartConnHealthy() {
		cartMu.Lock()
		defer cartMu.Unlock()
		if cartConn != nil {
			_ = cartConn.Close()
		}
		if cartClient == nil {
			ip, port := utils.Config.Cart.IP, utils.Config.Cart.Port
			var err error
			cartConn, err = grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return nil, err
			}
			log.Printf("已建立购物车服务的gRPC连接")
			cartClient = cart.NewCartServiceClient(cartConn)
		}
	}

	return cartClient, nil
}

func isCartConnHealthy() bool {
	state := cartConn.GetState()
	return state == connectivity.Ready || state == connectivity.Idle
}
