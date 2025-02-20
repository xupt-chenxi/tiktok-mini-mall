package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"tiktok-mini-mall/api/pb/shop"
	"tiktok-mini-mall/internal/app/pkg/grpcclient"
)

type PlaceOrderReq struct {
	UserId     string  `json:"userId"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Address    string  `json:"address"`
	Amount     float32 `json:"amount"`
	OrderItems string  `json:"orderItems"`
}

type UpdateOrderStateReq struct {
	UserId  string `json:"userId"`
	OrderId string `json:"orderId"`
	State   uint32 `json:"state"`
}

func PlaceOrderHandler(c *gin.Context) {
	traceID := c.GetString("TraceID")
	var req PlaceOrderReq
	_ = c.ShouldBindJSON(&req)

	md := metadata.Pairs("trace-id", traceID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client, err := grpcclient.GetShopClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与shop服务建立连接失败: %v\n", traceID, err)
		return
	}

	var orderItems []*shop.OrderItem
	_ = json.Unmarshal([]byte(req.OrderItems), &orderItems)
	res, err := client.PlaceOrder(ctx, &shop.PlaceOrderReq{
		UserId:     req.UserId,
		Name:       req.Name,
		Email:      req.Email,
		Address:    req.Address,
		Amount:     req.Amount,
		OrderItems: orderItems,
	})
	if err != nil {
		sts := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": sts.Message()})
		log.Printf("TraceID: %v, 调用shop服务PlaceOrder返回错误: %v\n", traceID, sts.Message())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": res,
	})
}

func ListOrderHandler(c *gin.Context) {
	traceID := c.GetString("TraceID")
	userId := c.PostForm("userId")

	md := metadata.Pairs("trace-id", traceID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client, err := grpcclient.GetShopClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与shop服务建立连接失败: %v\n", traceID, err)
		return
	}
	res, err := client.ListOrder(ctx, &shop.ListOrderReq{
		UserId: userId,
	})

	if err != nil {
		sts := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": sts.Message()})
		log.Printf("TraceID: %v, 调用shop服务ListOrder返回错误: %v\n", traceID, sts.Message())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": res,
	})
}

func UpdateOrderState(c *gin.Context) {
	traceID := c.GetString("TraceID")
	var req UpdateOrderStateReq
	_ = c.ShouldBindJSON(&req)
	md := metadata.Pairs("trace-id", traceID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client, err := grpcclient.GetShopClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与shop服务建立连接失败: %v\n", traceID, err)
		return
	}
	res, err := client.UpdateOrderState(ctx, &shop.UpdateOrderStateReq{
		UserId:  req.UserId,
		OrderId: req.OrderId,
		State:   req.State,
	})

	if err != nil {
		sts := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": sts.Message()})
		log.Printf("TraceID: %v, 调用shop服务UpdateOrderState返回错误: %v\n", traceID, sts.Message())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": res,
	})
}
