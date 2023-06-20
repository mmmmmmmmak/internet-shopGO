package product_usecase

import (
	"context"
	db_dto "main/internal/adapters/dto"
	"main/internal/apperror"
	"main/internal/domain/entity"
)

type Service interface {
	GetAll(ctx context.Context) ([]entity.Product, error)
	Add(ctx context.Context, dto db_dto.CreateProductDTO) (string, error)
}

type TokenManager interface {
	ValidateToken(tokenString string) (bool, error)
	TokenUser(accessToken string) (string, error)
}

type productUsecase struct {
	productService Service
	tokenManager   TokenManager
}

func (p productUsecase) GetAll(ctx context.Context) ([]entity.Product, error) {
	products, err := p.productService.GetAll(ctx)
	if err != nil {
		return nil, apperror.NewAppError(err, "can`t get products", "", "PR-000001")
	}
	return products, nil
}

func (p productUsecase) AddProduct(ctx context.Context, dto CreateProductDTO) (string, error) {
	sellerID, err := p.tokenManager.TokenUser(dto.Token)
	if err != nil {
		return "", apperror.NewAppError(err, "can`t add product", "", "PR-000003")
	}
	dbdto := db_dto.CreateProductDTO{
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		Seller:      sellerID,
	}
	idProduct, err := p.productService.Add(ctx, dbdto)
	if err != nil {
		return "", apperror.NewAppError(err, "can`t add product", "", "PR-000003")
	}
	return idProduct, err
}

// создать функцию создания товара (нужно к юзеру добавить массив его товаров, статус селлер и покупатель)
func NewProductUsecase(service Service, tokenManager TokenManager) *productUsecase {
	return &productUsecase{
		productService: service,
		tokenManager:   tokenManager,
	}
}
