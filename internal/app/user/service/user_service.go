// Package service 提供用户服务
// Author: chenxi 2025.01
package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"tiktok-mini-mall/api/pb/user"
	"tiktok-mini-mall/internal/app/user/errortype"
	"tiktok-mini-mall/internal/app/user/model"
	"tiktok-mini-mall/internal/app/user/repository"
	"tiktok-mini-mall/pkg/utils"
	"time"
)

type UserService struct {
	user.UnimplementedUserServiceServer
}

// Register 注册功能
func (UserService) Register(ctx context.Context, req *user.RegisterReq) (*user.RegisterResp, error) {
	email, pass, confirmPass := req.GetEmail(), req.GetPassword(), req.GetConfirmPassword()
	md, _ := metadata.FromIncomingContext(ctx)
	traceID := md["trace-id"]
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
		err = errors.Wrap(err, "snowflake.NewNode 出错")
		log.Printf("TraceID: %v, err: %+v", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	// 生成雪花 ID
	snowID := node.Generate()
	err = repository.CreateUser(&model.User{
		Id:       snowID.Int64(),
		Email:    email,
		PassHash: hashPassword(req.GetPassword()),
		Nickname: genNickname(),
	})
	if err != nil {
		err = errors.Wrap(err, "创建用户失败")
		log.Printf("TraceID: %v, err: %+v", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &user.RegisterResp{
		UserId: snowID.String(),
	}, nil
}

// Login 登录功能
func (UserService) Login(ctx context.Context, req *user.LoginReq) (*user.LoginResp, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	traceID := md["trace-id"]

	email, pass := req.GetEmail(), req.GetPassword()
	userInfo, err := repository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.InvalidArgument, errortype.ErrUserNotFound.Error())
		}
		err = errors.Wrap(err, "查询用户出错")
		log.Printf("TraceID: %v, err: %+v", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	if hashPassword(pass) != userInfo.PassHash {
		return nil, status.Error(codes.InvalidArgument, errortype.ErrInvalidPassword.Error())
	}
	token := strings.ReplaceAll(uuid.New().String(), "-", "")
	ip := viper.GetString("redis.ip")
	port := viper.GetString("redis.port")
	password := viper.GetString("redis.password")
	dbStr := viper.GetString("redis.db")
	db, _ := strconv.Atoi(dbStr)
	redisClient := utils.NewRedisClient(ip+port, password, db)
	_ = redisClient.Set(context.Background(), "token:"+token, strconv.FormatInt(userInfo.Id, 10), 12*time.Hour)
	return &user.LoginResp{
		UserId:   strconv.FormatInt(userInfo.Id, 10),
		Token:    token,
		Nickname: userInfo.Nickname,
	}, nil
}

// 输入校验
func validateInput(email, pass, confirmPass string) error {
	if email == "" || pass == "" || confirmPass == "" {
		return errortype.ErrInputEmpty
	}
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errortype.ErrEmailFormat
	}
	if pass != confirmPass {
		return errortype.PasswordMismatch
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

	// 字符集，包含大小写字母和数字
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	nickname := ""

	// 生成 4 位字符
	for i := 0; i < 4; i++ {
		randomIndex := r.Intn(len(characters))      // 获取随机字符的索引
		nickname += string(characters[randomIndex]) // 累加字符
	}

	return "user_" + nickname
}
