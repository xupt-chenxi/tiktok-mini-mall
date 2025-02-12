package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"tiktok-mini-mall/api/pb/user"
	"tiktok-mini-mall/internal/app/pkg/grpcclient"
)

func RegisterHandler(c *gin.Context) {
	traceID := c.GetString("TraceID")
	var req user.RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),
			"code": http.StatusBadRequest})
		return
	}
	// 通过 gRPC 调用用户服务
	md := metadata.Pairs("trace-id", traceID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client, err := grpcclient.GetUserClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与用户服务建立连接失败: %v\n", traceID, err)
		return
	}
	res, err := client.Register(ctx, &req)

	if err != nil {
		sts := status.Convert(err)
		if sts.Code() == codes.Internal {
			c.JSON(http.StatusInternalServerError, gin.H{"error": sts.Message()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": sts.Message(),
				"code": http.StatusBadRequest})
		}
		log.Printf("TraceID: %v, 调用用户服务注册功能返回错误: %v\n", traceID, sts.Message())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": res,
	})
}

func LoginHandler(c *gin.Context) {
	traceID := c.GetString("TraceID")
	var req user.LoginReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),
			"code": http.StatusBadRequest})
		return
	}
	// 通过 gRPC 调用用户服务
	md := metadata.Pairs("trace-id", traceID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client, err := grpcclient.GetUserClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与用户服务建立连接失败: %v\n", traceID, err)
		return
	}
	res, err := client.Login(ctx, &req)

	if err != nil {
		sts := status.Convert(err)
		if sts.Code() == codes.Internal {
			c.JSON(http.StatusInternalServerError, gin.H{"error": sts.Message()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": sts.Message(),
				"code": http.StatusBadRequest})
		}
		log.Printf("TraceID: %v, 调用用户服务登录功能返回错误: %v\n", traceID, sts.Message())
		return
	}
	log.Printf("userInfo: %v", res)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": res,
	})
}
