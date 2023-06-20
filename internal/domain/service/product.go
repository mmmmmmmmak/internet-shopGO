package service

import (
	"context"
	db_dto "main/internal/adapters/dto"
	"main/internal/domain/entity"
)

type ProductStorage interface {
	GetProducts(ctx context.Context) ([]entity.Product, error)
	Add(ctx context.Context, dto db_dto.CreateProductDTO) (string, error)
}

type productService struct {
	storage ProductStorage
}

func (p productService) GetAll(ctx context.Context) ([]entity.Product, error) {
	return p.storage.GetProducts(ctx)
}

func (p productService) Add(ctx context.Context, dto db_dto.CreateProductDTO) (string, error) {
	return p.storage.Add(ctx, dto)
}

func NewProductService(storage ProductStorage) *productService {
	return &productService{storage: storage}
}
