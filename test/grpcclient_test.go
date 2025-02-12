package test

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"tiktok-mini-mall/api/pb/user"
	"tiktok-mini-mall/internal/app/pkg/grpcclient"
)

// old: 每次调用服务都创建一次 gRPC 连接
func BenchmarkOldGRPCClient(b *testing.B) {
	ip, port := "localhost", ":50051"

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn, _ := grpc.NewClient(ip+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
			client := user.NewUserServiceClient(conn)

			_, _ = client.Login(context.Background(), &user.LoginReq{
				Email:    "123456@example.com",
				Password: "123456",
			})
			_ = conn.Close()
		}
	})
}

// new: 复用同一条 gRPC 连接
func BenchmarkNewGRPCClient(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			client, _ := grpcclient.GetUserClient()

			_, _ = client.Login(context.Background(), &user.LoginReq{
				Email:    "123456@example.com",
				Password: "123456",
			})
		}
	})
}
