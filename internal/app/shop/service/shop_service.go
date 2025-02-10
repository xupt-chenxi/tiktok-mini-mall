// Package service shop 服务
// Author: chenxi 2025.02
package service

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
	"tiktok-mini-mall/api/pb/prod"
	"tiktok-mini-mall/api/pb/shop"
	"tiktok-mini-mall/internal/app/shop/model"
	"tiktok-mini-mall/internal/app/shop/repository"
	"tiktok-mini-mall/pkg/utils"
)

type ShopService struct {
	shop.UnimplementedShopServiceServer
}

func (ShopService) PlaceOrder(ctx context.Context, req *shop.PlaceOrderReq) (*shop.PlaceOrderResp, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	traceID := md["trace-id"]

	// 获取商品服务
	ip, port := utils.Config.Product.IP, utils.Config.Product.Port
	conn, err := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)
	client := prod.NewProductCatalogServiceClient(conn)
	orderItems := req.GetOrderItems()
	for _, orderItem := range orderItems {
		// 调用商品服务扣减库存
		_, err := client.DecreaseStock(ctx, &prod.DecreaseStockReq{
			Id:       orderItem.ProductId,
			Quantity: orderItem.Quantity,
		})
		if err != nil {
			sts := status.Convert(err)
			log.Printf("TraceID: %v, 调用shop服务DecreaseStock返回错误: %v\n", traceID, sts.Message())
			return nil, err
		}
	}

	node, err := snowflake.NewNode(1)
	if err != nil {
		err = errors.Wrap(err, "snowflake.NewNode 出错")
		log.Printf("TraceID: %v, err: %+v", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	// 生成雪花 ID
	snowID := node.Generate()

	orderItemsStr, _ := json.Marshal(req.OrderItems)
	userId, _ := strconv.ParseInt(req.GetUserId(), 10, 64)
	err = repository.AddOrder(&model.Order{
		Id:         snowID.Int64(),
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

func (ShopService) MarkOrderPaid(ctx context.Context, req *shop.MarkOrderPaidReq) (*shop.MarkOrderPaidResp, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	traceID := md["trace-id"]

	orderId, _ := strconv.ParseInt(req.GetOrderId(), 10, 64)
	userId, _ := strconv.ParseInt(req.GetUserId(), 10, 64)
	err := repository.MarkOrderPaid(userId, orderId)
	if err != nil {
		err = errors.Wrap(err, "订单支付失败")
		log.Printf("TraceID: %v, err: %+v", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}
