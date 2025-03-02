// Package service shop 服务
// Author: chenxi 2025.02
package service

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-clients/golang"
	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
	"tiktok-mini-mall/api/pb/prod"
	"tiktok-mini-mall/api/pb/shop"
	"tiktok-mini-mall/internal/app/pkg/grpcclient"
	"tiktok-mini-mall/internal/app/shop/model"
	"tiktok-mini-mall/internal/app/shop/repository"
	"tiktok-mini-mall/pkg/utils"
	"time"
)

type ShopService struct {
	shop.UnimplementedShopServiceServer
}

func (ShopService) PlaceOrder(ctx context.Context, req *shop.PlaceOrderReq) (*shop.PlaceOrderResp, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	traceID := md["trace-id"]

	// 获取商品服务
	client, err := grpcclient.GetProdClient()
	if err != nil {
		log.Printf("TraceID: %v, 与商品服务建立连接失败: %v\n", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	node, err := snowflake.NewNode(1)
	if err != nil {
		err = errors.Wrap(err, "snowflake.NewNode 出错")
		log.Printf("TraceID: %v, err: %+v", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	// 生成雪花 ID
	snowID := node.Generate()
	orderId := snowID.Int64()

	orderItemsStr, _ := json.Marshal(req.OrderItems)
	orderItems := req.GetOrderItems()
	var prodOrderItems []*prod.DecreaseStockItem
	for _, orderItem := range orderItems {
		prodOrderItems = append(prodOrderItems, &prod.DecreaseStockItem{
			ProductId: orderItem.ProductId,
			Quantity:  orderItem.Quantity,
		})
	}
	// 调用商品服务扣减库存
	_, err = client.DecreaseStock(ctx, &prod.DecreaseStockReq{
		OrderItems: prodOrderItems,
		OrderId:    strconv.Itoa(int(orderId)),
		UserId:     req.UserId,
	})
	if err != nil {
		sts := status.Convert(err)
		log.Printf("TraceID: %v, 调用shop服务DecreaseStock返回错误: %v\n", traceID, sts.Message())
		return nil, err
	}

	userId, _ := strconv.ParseInt(req.GetUserId(), 10, 64)
	err = repository.AddOrder(&model.Order{
		Id:         orderId,
		UserId:     userId,
		Name:       req.GetName(),
		Email:      req.GetEmail(),
		Address:    req.GetAddress(),
		Price:      req.GetAmount(),
		OrderItems: string(orderItemsStr),
		State:      0,
	})

	if err != nil {
		err = errors.Wrap(err, "生成订单失败")
		log.Printf("TraceID: %v, err: %+v", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 基于RocketMQ定时取消订单
	topic := utils.Config.RocketMQ.TopicShop
	producer, err := utils.NewProducer(topic)
	if err != nil {
		err = errors.Wrap(err, "RocketMQ新建生产者实例出错")
		log.Printf("TraceID: %v, err: %+v\n", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	cancelOrder := &struct {
		OrderId int64 `json:"orderId"`
		UserId  int64 `json:"userId"`
	}{
		OrderId: orderId,
		UserId:  userId,
	}
	data, _ := json.Marshal(cancelOrder)
	msg := &golang.Message{
		Topic: topic,
		Body:  data,
	}
	msg.SetKeys(snowID.String())
	msg.SetTag("tag_order")
	msg.SetDelayTimestamp(time.Now().Add(time.Second * 10))
	_, err = producer.Send(context.TODO(), msg)
	if err != nil {
		err = errors.Wrap(err, "向RocketMQ中发送订单定时取消信息出错")
		log.Printf("TraceID: %v, err: %+v\n", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	log.Printf("向RocketMQ中发送订单定时取消信息成功")

	return &shop.PlaceOrderResp{
		OrderId: snowID.String(),
	}, nil
}

func (ShopService) ListOrder(ctx context.Context, req *shop.ListOrderReq) (*shop.ListOrderResp, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	traceID := md["trace-id"]

	userId, _ := strconv.ParseInt(req.GetUserId(), 10, 64)
	orderList, err := repository.GetListOrder(userId)
	if err != nil {
		err = errors.Wrap(err, "获取订单列表失败")
		log.Printf("TraceID: %v, err: %+v", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	var orders []*shop.Order
	for _, order := range orderList {
		var orderItems []*shop.OrderItem
		_ = json.Unmarshal([]byte(order.OrderItems), &orderItems)
		orders = append(orders, &shop.Order{
			OrderItems: orderItems,
			OrderId:    strconv.FormatInt(order.Id, 10),
			State:      uint32(order.State),
		})
	}
	return &shop.ListOrderResp{
		Orders: orders,
	}, nil
}

func (ShopService) UpdateOrderState(ctx context.Context, req *shop.UpdateOrderStateReq) (*shop.UpdateOrderStateResp, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	traceID := md["trace-id"]

	orderId, _ := strconv.ParseInt(req.GetOrderId(), 10, 64)
	userId, _ := strconv.ParseInt(req.GetUserId(), 10, 64)
	err := repository.UpdateOrderState(userId, orderId, uint8(req.State))
	if err != nil {
		err = errors.Wrap(err, "订单支付失败")
		log.Printf("TraceID: %v, err: %+v", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}
