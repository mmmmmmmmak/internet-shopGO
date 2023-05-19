package service

import (
	"context"
	"main/internal/adapters/dto"
	"main/internal/domain/entity"
)

type UserStorage interface {
	Create(ctx context.Context, user db_dto.CreateUserDTO) (string, error)
	IsUserCreated(ctx context.Context, user db_dto.IsUserExists) (bool, error)
	AuthUser(ctx context.Context, user db_dto.AuthUser) (string, error)
	Update(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, id string) error
}

type userService struct {
	storage UserStorage
}

func NewUserService(storage UserStorage) *userService {
	return &userService{storage: storage}
}

func (s userService) Create(ctx context.Context, dbDTO db_dto.CreateUserDTO) (string, error) {
	return s.storage.Create(ctx, dbDTO)
}

func (s userService) IsUserCreated(ctx context.Context, user db_dto.IsUserExists) (bool, error) {
	return s.storage.IsUserCreated(ctx, user)
}

func (s userService) AuthUser(ctx context.Context, user db_dto.AuthUser) (string, error) {
	return s.storage.AuthUser(ctx, user)
}
