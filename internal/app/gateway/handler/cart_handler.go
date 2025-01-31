package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"strconv"
	"tiktok-mini-mall/api/pb/cart_pb"
)

type AddItemReq struct {
	UserId    int64  `json:"userId"`
	ProductId uint32 `json:"productId"`
	Quantity  int32  `json:"quantity"`
}

func AddItemHandler(c *gin.Context) {
	traceID := c.GetString("TraceID")
	var req AddItemReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("ShouldBindJSON err: %v\n", err)
	}

	md := metadata.Pairs("trace-id", traceID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	ip, port := viper.GetString("cart.ip"), viper.GetString("cart.port")
	conn, err := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与cart服务建立连接失败: %v\n", traceID, err)
		return
	}
	defer conn.Close()
	client := cart.NewCartServiceClient(conn)
	res, err := client.AddItem(ctx, &cart.AddItemReq{
		UserId: req.UserId,
		Item: &cart.CartItem{
			ProductId: req.ProductId,
			Quantity:  req.Quantity,
		},
	})
	if err != nil {
		sts := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": sts.Message()})
		log.Printf("TraceID: %v, 调用cart服务AddItem返回错误: %v\n", traceID, sts.Message())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": res,
	})
}

func GetCartHandler(c *gin.Context) {
	traceID := c.GetString("TraceID")
	userIdParm := c.PostForm("userId")

	md := metadata.Pairs("trace-id", traceID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	ip, port := viper.GetString("cart.ip"), viper.GetString("cart.port")
	conn, err := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与cart服务建立连接失败: %v\n", traceID, err)
		return
	}
	defer conn.Close()
	client := cart.NewCartServiceClient(conn)
	userId, err := strconv.Atoi(userIdParm)
	res, err := client.GetCart(ctx, &cart.GetCartReq{
		UserId: int64(userId),
	})
	if err != nil {
		sts := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": sts.Message()})
		log.Printf("TraceID: %v, 调用cart服务GetCart返回错误: %v\n", traceID, sts.Message())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": res,
	})
}

func EmptyCartHandler(c *gin.Context) {
	traceID := c.GetString("TraceID")
	userIdParm := c.PostForm("userId")

	md := metadata.Pairs("trace-id", traceID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	ip, port := viper.GetString("cart.ip"), viper.GetString("cart.port")
	conn, err := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与cart服务建立连接失败: %v\n", traceID, err)
		return
	}
	defer conn.Close()
	client := cart.NewCartServiceClient(conn)
	userId, err := strconv.Atoi(userIdParm)
	res, err := client.EmptyCart(ctx, &cart.EmptyCartReq{
		UserId: int64(userId),
	})
	if err != nil {
		sts := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": sts.Message()})
		log.Printf("TraceID: %v, 调用cart服务EmptyCart返回错误: %v\n", traceID, sts.Message())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": res,
	})
}
