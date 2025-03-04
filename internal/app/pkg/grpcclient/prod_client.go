package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"tiktok-mini-mall/api/pb/prod"
	"tiktok-mini-mall/pkg/utils"
)

var (
	prodClient prod.ProductCatalogServiceClient
	prodConn   *grpc.ClientConn
	prodMu     sync.Mutex
)

func GetProdClient() (prod.ProductCatalogServiceClient, error) {
	if prodClient == nil || !isProdConnHealthy() {
		prodMu.Lock()
		defer prodMu.Unlock()
		if prodConn != nil {
			_ = prodConn.Close()
		}
		if prodClient == nil || !isProdConnHealthy() {
			ip, port := utils.Config.Product.IP, utils.Config.Product.Port
			var err error
			prodConn, err = grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return nil, err
			}
			log.Printf("已建立商品服务的gRPC连接")
			prodClient = prod.NewProductCatalogServiceClient(prodConn)
		}
	}

	return prodClient, nil
}

func isProdConnHealthy() bool {
	state := prodConn.GetState()
	return state == connectivity.Ready || state == connectivity.Idle
}
