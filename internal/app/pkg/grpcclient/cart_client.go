package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"tiktok-mini-mall/api/pb/cart"
	"tiktok-mini-mall/pkg/utils"
)

var cartClient cart.CartServiceClient

// var cartConn *grpc.ClientConn
var cartMu sync.Mutex

func GetCartClient() (cart.CartServiceClient, error) {
	if cartClient == nil {
		cartMu.Lock()
		defer cartMu.Unlock()
		if cartClient == nil {
			ip, port := utils.Config.Cart.IP, utils.Config.Cart.Port
			cartConn, err := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return nil, err
			}
			log.Printf("已建立购物车服务的gRPC连接")
			cartClient = cart.NewCartServiceClient(cartConn)
		}
	}

	return cartClient, nil
}

//func CloseCartConn() error {
//	if cartConn != nil {
//		err := cartConn.Close()
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
