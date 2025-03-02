package main

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-clients/golang"
	"github.com/apache/rocketmq-clients/golang/credentials"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"tiktok-mini-mall/api/pb/shop"
	"tiktok-mini-mall/internal/app/pkg/grpcclient"
	"tiktok-mini-mall/internal/app/shop/repository"
	"tiktok-mini-mall/internal/app/shop/service"
	"tiktok-mini-mall/pkg/utils"
	"time"
)

func main() {
	dsn := utils.Config.Shop.Database.DSN
	repository.InitDatabase(dsn)

	namingClient, err := utils.NewNamingClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          utils.Config.Shop.IP,
		Port:        50054,
		ServiceName: "shop_service",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	if success == false {
		log.Fatalf("shop服务注册失败")
	}
	log.Println("shop服务注册成功")
	go processCancelOrder()
	port := utils.Config.Shop.Port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("shop服务监听端口失败: %v", err)
	}
	grpcServer := grpc.NewServer()
	shop.RegisterShopServiceServer(grpcServer, &service.ShopService{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("shop服务启动失败: %v", err)
	}
}

func processCancelOrder() {
	client, err := grpcclient.GetShopClient()
	if err != nil {
		log.Fatalf("与shop服务建立连接失败: %v\n", err)
	}
	simpleConsumer, err := golang.NewSimpleConsumer(&golang.Config{
		Endpoint:      utils.Config.RocketMQ.Endpoint,
		ConsumerGroup: "group_order",
		Credentials: &credentials.SessionCredentials{
			AccessKey:    utils.Config.RocketMQ.AccessKey,
			AccessSecret: utils.Config.RocketMQ.SecretKey,
		},
	},
		golang.WithAwaitDuration(time.Second*5),
		golang.WithSubscriptionExpressions(map[string]*golang.FilterExpression{
			utils.Config.RocketMQ.TopicShop: golang.SUB_ALL,
		}),
	)
	if err != nil {
		log.Printf("拉取RocketMQ订单定时取消失败: %v\n", err)
	}
	err = simpleConsumer.Start()
	if err != nil {
		log.Printf("拉取RocketMQ订单定时取消失败: %v\n", err)
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
			log.Printf("拉取RocketMQ订单定时取消失败: %v\n", err)
		}
		for _, mv := range mvs {
			err = simpleConsumer.Ack(context.TODO(), mv)
			if err != nil {
				log.Printf("拉取RocketMQ订单定时取消失败: %v\n", err)
			}
			cancelOrder := &struct {
				OrderId int64 `json:"orderId"`
				UserId  int64 `json:"userId"`
			}{}
			_ = json.Unmarshal(mv.GetBody(), &cancelOrder)
			orderId := cancelOrder.OrderId
			state, err := repository.GetOrderState(orderId)
			if err != nil {
				log.Printf("获取订单状态失败: %v\n", err)
			}
			if state == 1 {
				_, err := client.UpdateOrderState(context.Background(), &shop.UpdateOrderStateReq{
					UserId:  strconv.FormatInt(cancelOrder.UserId, 10),
					OrderId: strconv.FormatInt(orderId, 10),
					State:   4,
				})
				if err != nil {
					log.Printf("对于订单号: %v 状态更新失败: %v", orderId, err)
				} else {
					log.Printf("订单号: %v 已超时取消", orderId)
				}
			}
		}
	}
}
