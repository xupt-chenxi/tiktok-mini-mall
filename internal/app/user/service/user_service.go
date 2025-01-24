// 提供用户服务
// Author: chenxi 2025.01
package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"math/rand"
	"regexp"
	"tiktok-mini-mall/api/pb/user_pb"
	"tiktok-mini-mall/internal/app/user/model"
	"tiktok-mini-mall/internal/app/user/repository"
	"time"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
}

// Register 注册功能
func (UserService) Register(ctx context.Context, req *userpb.RegisterReq) (*userpb.RegisterResp, error) {
	email, pass, confirmPass := req.GetEmail(), req.GetPassword(), req.GetConfirmPassword()
	// 1.输入校验
	err := validateInput(email, pass, confirmPass)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// 2.检测用户是否已存在
	_, err = repository.GetUserByEmail(email)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "用户已存在，无需再进行注册")
	}
	// 3.用户注册
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	// 生成雪花 ID
	snowID := node.Generate()
	err = repository.CreateUser(&model.User{
		UserID:   snowID.Int64(),
		Email:    email,
		PassHash: hashPassword(req.GetPassword()),
		Nickname: genNickname(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	log.Printf("用户注册成功, 用户 id 为: %d", snowID.Int64())
	return &userpb.RegisterResp{
		UserId: snowID.Int64(),
	}, nil
}

// Login 登录功能
func (UserService) Login(ctx context.Context, req *userpb.LoginReq) (*userpb.LoginResp, error) {
	email, pass := req.GetEmail(), req.GetPassword()
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if hashPassword(pass) != user.PassHash {
		return nil, status.Error(codes.InvalidArgument, errors.New("密码错误").Error())
	}

	return &userpb.LoginResp{
		UserId: user.UserID,
	}, nil
}

// 输入校验
func validateInput(email, pass, confirmPass string) error {
	if email == "" || pass == "" || confirmPass == "" {
		return errors.New("用户输入不能有空")
	}
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errors.New("用户输入邮箱格式有误")
	}
	if pass != confirmPass {
		return errors.New("用户两次输入的密码不一致")
	}

	return nil
}

// 对密码进行哈希
func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// 生成随机昵称
func genNickname() string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	timestamp := time.Now().UnixNano()
	randomSuffix := r.Intn(1000)
	return fmt.Sprintf("User_%d_%d", timestamp, randomSuffix)
}
