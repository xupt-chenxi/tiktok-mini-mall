package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"tiktok-mini-mall/api/pb/prod"
	"tiktok-mini-mall/internal/app/product/repository"
	"tiktok-mini-mall/internal/app/product/service"
	"tiktok-mini-mall/pkg/utils"
)

func main() {
	dsn := utils.Config.Product.Database.DSN
	repository.InitDatabase(dsn)

	port := utils.Config.Product.Port
	// 注册商品服务
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("商品服务监听端口失败: %v", err)
	}
	grpcServer := grpc.NewServer()
	prod.RegisterProductCatalogServiceServer(grpcServer, &service.ProductService{})
	err = cacheStock()
	if err != nil {
		log.Fatal("商品库存缓存预热失败: ", err)
	}
	log.Println("商品服务启动...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("商品服务启动失败: %v", err)
	}
}

// 库存缓存预热
func cacheStock() error {
	stockList, err := repository.GetStockList()
	if err != nil {
		return err
	}
	ip, port, password, dbStr := utils.Config.Redis.IP, utils.Config.Redis.Port, utils.Config.Redis.Password, utils.Config.Redis.DB
	db, _ := strconv.Atoi(dbStr)
	redisClient := utils.NewRedisClient(ip+port, password, db)
	for _, value := range stockList {
		productId, stock := value.Id, value.Stock
		prodIdStr, stockStr := strconv.Itoa(int(productId)), strconv.Itoa(int(stock))
		err := redisClient.Set(context.Background(), "stock:"+prodIdStr, stockStr, 0)
		if err != nil {
			return err
		}
	}
	return nil
}
