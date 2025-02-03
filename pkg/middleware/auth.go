package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"tiktok-mini-mall/pkg/utils"
)

const UserUnLogin = "用户未登录"

func AuthMiddleware() gin.HandlerFunc {
	// 不需要拦截的路径
	whitelist := map[string]bool{
		"/user/login":            true, // 登录接口
		"/user/register":         true, // 注册接口
		"/product/list-products": true, // 商品列表接口
	}

	return func(c *gin.Context) {
		if _, ok := whitelist[c.Request.URL.Path]; ok {
			c.Next()
			return
		}

		token := c.GetHeader("Authorization")
		// token 遵守 Bearer 规范
		if len(token) < 7 || token[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": UserUnLogin})
			return
		}
		token = token[7:]

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  http.StatusUnauthorized,
				"error": UserUnLogin,
			})
			return
		}
		utils.InitViper("configs/config.yaml")
		ip := viper.GetString("redis.ip")
		port := viper.GetString("redis.port")
		password := viper.GetString("redis.password")
		dbStr := viper.GetString("redis.db")
		db, _ := strconv.Atoi(dbStr)
		redisClient := utils.NewRedisClient(ip+port, password, db)
		userId, _ := redisClient.Get(context.Background(), "token:"+token)
		if userId == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  http.StatusUnauthorized,
				"error": UserUnLogin,
			})
			return
		}

		c.Next()
	}
}
