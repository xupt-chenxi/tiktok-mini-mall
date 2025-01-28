// 提供商品服务
// Author: chenxi 2025.01
package service

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	prod "tiktok-mini-mall/api/pb/prod_pb"
	"tiktok-mini-mall/internal/app/product/repository"
)

type ProductService struct {
	prod.UnimplementedProductCatalogServiceServer
}

func (ProductService) ListProducts(ctx context.Context, req *prod.ListProductsReq) (*prod.ListProductsResp, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	traceID := md["trace-id"]

	page, pageSize, categoryName := req.GetPage(), req.GetPageSize(), req.GetCategoryName()
	productList, err := repository.GetProductList(int(page), int(pageSize), categoryName)
	if err != nil {
		err = errors.Wrap(err, "查询商品列表出错")
		log.Printf("TraceID: %v, err: %+v\n", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	var products []*prod.Product
	for _, product := range productList {
		var categories []string
		_ = json.Unmarshal([]byte(product.Categories), &categories)
		products = append(products, &prod.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Picture:     product.Picture,
			Price:       product.Price,
			Categories:  categories,
		})
	}

	return &prod.ListProductsResp{
		Products: products,
	}, err
}

func (ProductService) GetProduct(ctx context.Context, req *prod.GetProductReq) (*prod.GetProductResp, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	traceID := md["trace-id"]

	product, err := repository.GetProductById(int(req.GetId()))
	if err != nil {
		err = errors.Wrap(err, "查询商品信息出错")
		log.Printf("TraceID: %v, err: %+v\n", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	var categories []string
	_ = json.Unmarshal([]byte(product.Categories), &categories)
	return &prod.GetProductResp{
		Product: &prod.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Picture:     product.Picture,
			Price:       product.Price,
			Categories:  categories,
		},
	}, nil
}

func (ProductService) SearchProducts(ctx context.Context, req *prod.SearchProductsReq) (*prod.SearchProductsResp, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	traceID := md["trace-id"]

	productList, err := repository.SearchProducts(req.GetQuery())
	if err != nil {
		err = errors.Wrap(err, "搜索商品出错")
		log.Printf("TraceID: %v, err: %+v\n", traceID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	var products []*prod.Product
	for _, product := range productList {
		var categories []string
		_ = json.Unmarshal([]byte(product.Categories), &categories)
		products = append(products, &prod.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Picture:     product.Picture,
			Price:       product.Price,
			Categories:  categories,
		})
	}
	return &prod.SearchProductsResp{
		Results: products,
	}, nil
}
