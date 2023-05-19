package service

import (
	"context"
	"main/internal/adapters/dto"
	"main/internal/domain/entity"
)

type UserStorage interface {
	Create(ctx context.Context, user db_dto.CreateUserDTO) (string, error)
	IsUserCreated(ctx context.Context, user db_dto.IsUserExists) (bool, error)
	FindByEmail(ctx context.Context, user db_dto.AuthByEmail) (string, error)
	FindByUsername(ctx context.Context, user db_dto.AuthByUsername) (string, error)
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
	userID, err := s.storage.Create(ctx, dbDTO)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (s userService) IsUserCreated(ctx context.Context, user db_dto.IsUserExists) (bool, error) {
	find, err := s.storage.IsUserCreated(ctx, user)
	if err != nil {
		return false, err
	}
	return find, nil
}

func (s userService) FindByEmail(ctx context.Context, user db_dto.AuthByEmail) (string, error) {
	id, err := s.storage.FindByEmail(ctx, user)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s userService) FindByUsername(ctx context.Context, user db_dto.AuthByUsername) (string, error) {
	id, err := s.storage.FindByUsername(ctx, user)
	if err != nil {
		return "", err
	}
	return id, nil
}
