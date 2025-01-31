package service

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	cart "tiktok-mini-mall/api/pb/cart_pb"
	"tiktok-mini-mall/internal/app/cart/repository"
)

type CartService struct {
	cart.UnimplementedCartServiceServer
}

func (CartService) AddItem(ctx context.Context, req *cart.AddItemReq) (*cart.AddItemResp, error) {
	userId, item := req.GetUserId(), req.GetItem()
	md, _ := metadata.FromIncomingContext(ctx)
	traceID := md["trace-id"]

	err := repository.AddItem(item, int(userId))
	if err != nil {
		err = errors.Wrap(err, "添加购物车失败")
		log.Printf("TraceID: %v, err: %+v", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

func (CartService) GetCart(ctx context.Context, req *cart.GetCartReq) (*cart.GetCartResp, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	traceID := md["trace-id"]

	userId := req.GetUserId()
	cartItems, err := repository.GetCart(int(userId))
	if err != nil {
		err = errors.Wrap(err, "查询购物车失败")
		log.Printf("TraceID: %v, err: %+v", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	var items []*cart.CartItem
	for _, item := range cartItems {
		items = append(items, &cart.CartItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}
	return &cart.GetCartResp{
		Cart: &cart.Cart{
			UserId: userId,
			Items:  items,
		},
	}, nil
}

func (CartService) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (*cart.EmptyCartResp, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	traceID := md["trace-id"]

	userId := req.GetUserId()
	err := repository.EmptyCart(int(userId))
	if err != nil {
		err = errors.Wrap(err, "清空购物车失败")
		log.Printf("TraceID: %v, err: %+v", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}
