package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	userpb "tiktok-mini-mall/api/pb/user_pb"
	"time"
)

func RegisterHandler(c *gin.Context) {
	var req userpb.RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),
			"code": http.StatusBadRequest})
		return
	}
	// 通过 gRPC 调用用户服务
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ip, port := viper.GetString("user.ip"), viper.GetString("user.port")
	conn, err := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("与用户服务建立连接失败: %v\n", err)
		return
	}
	defer conn.Close()
	client := userpb.NewUserServiceClient(conn)
	res, err := client.Register(ctx, &req)

	if err != nil {
		sts := status.Convert(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": sts.Message(),
			"code": http.StatusBadRequest})
		log.Printf("调用用户服务注册功能返回错误: %v\n", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": res,
	})
}

func LoginHandler(c *gin.Context) {
	var req userpb.LoginReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),
			"code": http.StatusBadRequest})
		return
	}
	// 通过 gRPC 调用用户服务
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ip, port := viper.GetString("user.ip"), viper.GetString("user.port")
	conn, err := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("与用户服务建立连接失败: %v\n", err)
		return
	}
	defer conn.Close()
	client := userpb.NewUserServiceClient(conn)
	res, err := client.Login(ctx, &req)

	if err != nil {
		sts := status.Convert(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": sts.Message(),
			"code": http.StatusBadRequest})
		log.Printf("调用用户服务登录功能返回错误: %v\n", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": res,
	})
}
