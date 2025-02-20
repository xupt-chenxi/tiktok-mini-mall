package main

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-clients/golang"
	"github.com/apache/rocketmq-clients/golang/credentials"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"tiktok-mini-mall/api/pb/prod"
	"tiktok-mini-mall/api/pb/shop"
	"tiktok-mini-mall/internal/app/pkg/grpcclient"
	"tiktok-mini-mall/internal/app/product/repository"
	"tiktok-mini-mall/internal/app/product/service"
	"tiktok-mini-mall/pkg/utils"
	"time"
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
	go processStockDecrease()
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

func processStockDecrease() {
	client, err := grpcclient.GetShopClient()
	if err != nil {
		log.Fatalf("与shop服务建立连接失败: %v\n", err)
	}
	simpleConsumer, err := golang.NewSimpleConsumer(&golang.Config{
		Endpoint:      utils.Config.RocketMQ.Endpoint,
		ConsumerGroup: "group_stock",
		Credentials: &credentials.SessionCredentials{
			AccessKey:    utils.Config.RocketMQ.AccessKey,
			AccessSecret: utils.Config.RocketMQ.SecretKey,
		},
	},
		golang.WithAwaitDuration(time.Second*5),
		golang.WithSubscriptionExpressions(map[string]*golang.FilterExpression{
			utils.Config.RocketMQ.TopicProd: golang.SUB_ALL,
		}),
	)
	if err != nil {
		log.Printf("拉取RocketMQ扣减库存失败: %v\n", err)
	}
	err = simpleConsumer.Start()
	if err != nil {
		log.Printf("拉取RocketMQ扣减库存失败: %v\n", err)
	}
	defer func(simpleConsumer golang.SimpleConsumer) {
		_ = simpleConsumer.GracefulStop()
	}(simpleConsumer)
	for {
		time.Sleep(time.Second * 1)

		mvs, err := simpleConsumer.Receive(context.TODO(), 16, time.Second*20)
		if err != nil {
			if err.Error() == "CODE: MESSAGE_NOT_FOUND, MESSAGE: no new message" {
				continue
			}
			log.Printf("拉取RocketMQ扣减库存失败: %v\n", err)
		}
		for _, mv := range mvs {
			err = simpleConsumer.Ack(context.TODO(), mv)
			if err != nil {
				log.Printf("拉取RocketMQ扣减库存失败: %v\n", err)
			}
			var stockDecrease *prod.DecreaseStockReq
			_ = json.Unmarshal(mv.GetBody(), &stockDecrease)
			orderItems := stockDecrease.OrderItems
			orderId := stockDecrease.OrderId
			flag := true
			for _, orderItem := range orderItems {
				err = repository.DecreaseStock(orderItem.ProductId, orderItem.Quantity)
				if err != nil {
					log.Printf("订单: %v, 商品: %v, 库存扣减: %v 失败", orderId, orderItem.ProductId, orderItem.Quantity)
					flag = false
					break
				} else {
					log.Printf("订单: %v, 商品: %v, 库存扣减: %v 成功", orderId, orderItem.ProductId, orderItem.Quantity)
				}
			}
			if flag {
				_, err := client.UpdateOrderState(context.Background(), &shop.UpdateOrderStateReq{
					UserId:  stockDecrease.UserId,
					OrderId: orderId,
					State:   1,
				})
				if err != nil {
					log.Printf("对于订单号: %v 状态更新失败: %v", orderId, err)
				}
			} else {
				_, err := client.UpdateOrderState(context.Background(), &shop.UpdateOrderStateReq{
					UserId:  stockDecrease.UserId,
					OrderId: orderId,
					State:   2,
				})
				if err != nil {
					log.Printf("对于订单号: %v 状态更新失败: %v", orderId, err)
				}
			}
		}
	}
}
