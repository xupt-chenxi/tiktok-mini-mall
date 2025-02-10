package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"strconv"
	"tiktok-mini-mall/api/pb/prod"
	"tiktok-mini-mall/pkg/utils"
)

func ListProductsHandler(c *gin.Context) {
	traceID := c.GetString("TraceID")
	category := c.DefaultQuery("category", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	md := metadata.Pairs("trace-id", traceID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	ip, port := utils.Config.Product.IP, utils.Config.Product.Port
	conn, err := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与商品服务建立连接失败: %v\n", traceID, err)
		return
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)
	client := prod.NewProductCatalogServiceClient(conn)
	res, err := client.ListProducts(ctx, &prod.ListProductsReq{
		Page:         int32(page),
		PageSize:     int64(pageSize),
		CategoryName: category,
	})
	if err != nil {
		sts := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": sts.Message()})
		log.Printf("TraceID: %v, 调用商品服务ListProducts返回错误: %v\n", traceID, sts.Message())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": res,
	})
}

func GetProductHandler(c *gin.Context) {
	traceID := c.GetString("TraceID")
	id, _ := strconv.Atoi(c.Param("id"))

	md := metadata.Pairs("trace-id", traceID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	ip, port := utils.Config.Product.IP, utils.Config.Product.Port
	conn, err := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与商品服务建立连接失败: %v\n", traceID, err)
		return
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)
	client := prod.NewProductCatalogServiceClient(conn)
	res, err := client.GetProduct(ctx, &prod.GetProductReq{
		Id: uint32(id),
	})
	if err != nil {
		sts := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": sts.Message()})
		log.Printf("TraceID: %v, 调用商品服务GetProduct返回错误: %v\n", traceID, sts.Message())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": res,
	})
}

func SearchProductsHandler(c *gin.Context) {
	traceID := c.GetString("TraceID")
	query := c.Query("query")

	md := metadata.Pairs("trace-id", traceID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	ip, port := utils.Config.Product.IP, utils.Config.Product.Port
	conn, err := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),
			"code": http.StatusInternalServerError})
		log.Printf("TraceID: %v, 与商品服务建立连接失败: %v\n", traceID, err)
		return
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)
	client := prod.NewProductCatalogServiceClient(conn)
	res, err := client.SearchProducts(ctx, &prod.SearchProductsReq{
		Query: query,
	})
	if err != nil {
		sts := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": sts.Message()})
		log.Printf("TraceID: %v, 调用商品服务SearchProducts返回错误: %v\n", traceID, sts.Message())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": res,
	})
}
