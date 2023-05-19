package user_usecase

import (
	"context"
	"github.com/asaskevich/govalidator"
	db_dto "main/internal/adapters/dto"
	"main/internal/apperror"
	"main/internal/config"
	"main/pkg/utils"
)

//это бизнес-сценарий моей системы

//скорее всего юзер будет не нужен, нужно будет для заказов (Orders)

type Service interface {
	Create(ctx context.Context, dto db_dto.CreateUserDTO) (string, error)
	IsUserCreated(ctx context.Context, user db_dto.IsUserExists) (bool, error)
	AuthUser(ctx context.Context, user db_dto.AuthUser) (string, error)
}

type userUsecase struct {
	userService Service
}

func (u userUsecase) CreateUser(ctx context.Context, dto CreateUserDTO) (string, error) {
	cfg := config.GetConfig()
	switch {
	case !govalidator.IsEmail(dto.Email):
		return "", apperror.NewAppError(nil, "email entered incorrectly", "", "US-000001")
	case dto.Password == "" || dto.Username == "":
		return "", apperror.NewAppError(nil, "one of the parameters is omitted", "", "US-000002")
	}

	user := db_dto.IsUserExists{
		Email:    dto.Email,
		Username: dto.Username,
	}
	isUserExists, err := u.userService.IsUserCreated(ctx, user)
	if err != nil {
		return "", err
	}
	if !isUserExists {
		return "", apperror.NewAppError(nil, "such a email or username already registered", "", "US-000006")
	}

	passwordHash := utils.HashPassword(dto.Password)

	dbDTO := db_dto.CreateUserDTO{
		Email:        dto.Email,
		Username:     dto.Username,
		PasswordHash: passwordHash,
	}
	userID, err := u.userService.Create(ctx, dbDTO)
	if err != nil {
		return "", apperror.NewAppError(err, "failed to create user", "", "US-000007")
	}
	jwtToken, err := utils.GenerateToken(userID, []byte(cfg.JWTConfig.Secret))
	if err != nil {
		return "", apperror.NewAppError(err, "unknown error", "failed to create jwt token", "US-000007")
	}
	return jwtToken, nil
}

func (u userUsecase) AuthUser(ctx context.Context, dto AuthUser) (string, error) {
	cfg := config.GetConfig()

	passwordHash := utils.HashPassword(dto.Password)
	userDTO := db_dto.AuthUser{
		Email:        dto.Email,
		Username:     dto.Username,
		PasswordHash: passwordHash,
	}
	userID, err := u.userService.AuthUser(ctx, userDTO)
	if err != nil {
		return "", err
	}

	jwtToken, err := utils.GenerateToken(userID, []byte(cfg.JWTConfig.Secret))
	if err != nil {
		return "", apperror.NewAppError(err, "unknown error", "failed to create jwt token", "US-000007")
	}
	return jwtToken, nil
}

func NewUserUsecase(service Service) *userUsecase {
	return &userUsecase{userService: service}
}

//func (u bookUsecase) ListAllBooks(ctx context.Context) []entity.BookView {
//	// отобразить список книг с именем Жанра и именем Автора
//	return u.bookService.GetAllForList(ctx)
//}
//
//func (u bookUsecase) GetFullBook(ctx context.Context, id string) entity.FullBook {
//	book := u.bookService.GetByID(ctx, id)
//	author := u.authorService.GetByID(ctx, book.AuthorID)
//	genre := u.genreService.GetByID(ctx, book.GenreID)
//
//	return entity.FullBook{
//		Book:   book,
//		Author: author,
//		Genre:  genre,
//	}
//}
//
//// pagination
//func (u bookUsecase) GetBooksWithAllAuthors(ctx context.Context, id string) []entity.BookView {
//	// Book{Authors: [all authors]}
//	// book, author(book_id) -=-
//	return nil
//}
