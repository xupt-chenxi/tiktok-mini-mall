package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"strconv"
	shop "tiktok-mini-mall/api/pb/shop_pb"
)

type PlaceOrderReq struct {
	UserId     int64   `json:"userId"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Address    string  `json:"address"`
	Amount     float32 `json:"amount"`
	OrderItems string  `json:"orderItems"`
}

type MarkOrderPaidReq struct {
	UserId  int64  `json:"userId"`
	OrderId string `json:"orderId"`
}

func PlaceOrderHandler(c *gin.Context) {
	traceID := c.GetString("TraceID")
	var req PlaceOrderReq
	_ = c.ShouldBindJSON(&req)

	md := metadata.Pairs("trace-id", traceID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	ip, port := viper.GetString("shop.ip"), viper.GetString("shop.port")
	conn, err := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与shop服务建立连接失败: %v\n", traceID, err)
		return
	}
	defer conn.Close()
	client := shop.NewShopServiceClient(conn)

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
	userIdStr := c.PostForm("userId")

	md := metadata.Pairs("trace-id", traceID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	ip, port := viper.GetString("shop.ip"), viper.GetString("shop.port")
	conn, err := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与shop服务建立连接失败: %v\n", traceID, err)
		return
	}
	defer conn.Close()
	client := shop.NewShopServiceClient(conn)
	userId, _ := strconv.Atoi(userIdStr)
	res, err := client.ListOrder(ctx, &shop.ListOrderReq{
		UserId: int64(userId),
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

func MarkOrderPaid(c *gin.Context) {
	traceID := c.GetString("TraceID")
	var req MarkOrderPaidReq
	_ = c.ShouldBindJSON(&req)
	md := metadata.Pairs("trace-id", traceID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	ip, port := viper.GetString("shop.ip"), viper.GetString("shop.port")
	conn, err := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与shop服务建立连接失败: %v\n", traceID, err)
		return
	}
	defer conn.Close()
	client := shop.NewShopServiceClient(conn)
	res, err := client.MarkOrderPaid(ctx, &shop.MarkOrderPaidReq{
		UserId:  req.UserId,
		OrderId: req.OrderId,
	})

	if err != nil {
		sts := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": sts.Message()})
		log.Printf("TraceID: %v, 调用shop服务MarkOrderPaid返回错误: %v\n", traceID, sts.Message())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": res,
	})
}
