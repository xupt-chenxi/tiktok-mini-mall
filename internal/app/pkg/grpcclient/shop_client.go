package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"tiktok-mini-mall/api/pb/shop"
	"tiktok-mini-mall/pkg/utils"
)

var shopClient shop.ShopServiceClient

// var shopConn *grpc.ClientConn
var shopMu sync.Mutex

func GetShopClient() (shop.ShopServiceClient, error) {
	if shopClient == nil {
		shopMu.Lock()
		defer shopMu.Unlock()
		if shopClient == nil {
			ip, port := utils.Config.Shop.IP, utils.Config.Shop.Port
			shopConn, err := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return nil, err
			}
			log.Printf("已建立Shop服务的gRPC连接")
			shopClient = shop.NewShopServiceClient(shopConn)
		}
	}

	return shopClient, nil
}

//func CloseShopConn() error {
//	if shopConn != nil {
//		err := shopConn.Close()
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
