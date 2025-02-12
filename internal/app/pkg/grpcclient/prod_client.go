package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"tiktok-mini-mall/api/pb/prod"
	"tiktok-mini-mall/pkg/utils"
)

var prodClient prod.ProductCatalogServiceClient

// var prodConn *grpc.ClientConn
var prodMu sync.Mutex

func GetProdClient() (prod.ProductCatalogServiceClient, error) {
	if prodClient == nil {
		prodMu.Lock()
		defer prodMu.Unlock()
		if prodClient == nil {
			ip, port := utils.Config.Product.IP, utils.Config.Product.Port
			prodConn, err := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return nil, err
			}
			log.Printf("已建立商品服务的gRPC连接")
			prodClient = prod.NewProductCatalogServiceClient(prodConn)
		}
	}

	return prodClient, nil
}

//func CloseProdConn() error {
//	if prodConn != nil {
//		err := prodConn.Close()
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
