package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"tiktok-mini-mall/api/pb/cart"
	"tiktok-mini-mall/internal/app/pkg/grpcclient"
)

type AddItemReq struct {
	UserId    string `json:"userId"`
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
	client, err := grpcclient.GetCartClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与cart服务建立连接失败: %v\n", traceID, err)
		return
	}
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
	userId := c.PostForm("userId")

	md := metadata.Pairs("trace-id", traceID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client, err := grpcclient.GetCartClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与cart服务建立连接失败: %v\n", traceID, err)
		return
	}
	res, err := client.GetCart(ctx, &cart.GetCartReq{
		UserId: userId,
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
	userId := c.PostForm("userId")

	md := metadata.Pairs("trace-id", traceID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client, err := grpcclient.GetCartClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与cart服务建立连接失败: %v\n", traceID, err)
		return
	}
	res, err := client.EmptyCart(ctx, &cart.EmptyCartReq{
		UserId: userId,
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
