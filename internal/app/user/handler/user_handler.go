package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	userpb "tiktok-mini-mall/api/pb/user_pb"
	"tiktok-mini-mall/internal/app/user/service"
)

var userService = service.UserService{}

func RegisterHandler(c *gin.Context) {
	var req userpb.RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := userService.Register(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func LoginHandler(c *gin.Context) {
	var req userpb.LoginReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := userService.Login(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
