package user_usecase

import (
	"context"
	"github.com/asaskevich/govalidator"
	db_dto "main/internal/adapters/dto"
	"main/internal/apperror"
	"main/internal/domain/entity"
	"main/pkg/utils"
)

//это бизнес-сценарий моей системы

//скорее всего юзер будет не нужен, нужно будет для заказов (Orders)

type Service interface {
	Create(ctx context.Context, dto db_dto.CreateUserDTO) (string, error)
	IsUserCreated(ctx context.Context, user db_dto.IsUserExistsDTO) (bool, error)
	AuthUser(ctx context.Context, user db_dto.AuthUserDTO) (string, error)
	GetUser(ctx context.Context, user db_dto.GetUserDTO) (entity.User, error)
	GetUserByRefreshToken(ctx context.Context, user db_dto.GetUserByRefreshTokenDTO) (string, error)
}

type TokenManager interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken() (string, error)
	TokenUser(accessToken string) (string, error)
	TokenExpires(accessToken string) (int64, error)
	ValidateToken(tokenString string) (bool, error)
}

type userUsecase struct {
	userService  Service
	tokenManager TokenManager
}

func (u userUsecase) CreateUser(ctx context.Context, dto CreateUserDTO) (Tokens, error) {
	var tokens Tokens
	switch {
	case !govalidator.IsEmail(dto.Email):
		return tokens, apperror.NewAppError(nil, "email entered incorrectly", "", "US-000001")
	case dto.Password == "" || dto.Username == "":
		return tokens, apperror.NewAppError(nil, "one of the parameters is omitted", "", "US-000002")
	}

	user := db_dto.IsUserExistsDTO{
		Email:    dto.Email,
		Username: dto.Username,
	}
	isUserExists, err := u.userService.IsUserCreated(ctx, user)
	if err != nil {
		return tokens, err
	}
	if !isUserExists {
		return tokens, apperror.NewAppError(nil, "such a email or username already registered", "", "US-000006")
	}

	passwordHash := utils.HashPassword(dto.Password)

	tokens.RefreshToken, err = u.tokenManager.GenerateRefreshToken()
	if err != nil {
		return tokens, apperror.NewAppError(err, "unknown error", "failed to create refresh jwt token", "TJ-000001")
	}

	expiresAt, err := u.tokenManager.TokenExpires(tokens.RefreshToken)
	if err != nil {
		return tokens, apperror.NewAppError(err, "token error", "the refresh token was compiled incorrectly", "TJ-000002")
	}

	dbDTO := db_dto.CreateUserDTO{
		Email:        dto.Email,
		Username:     dto.Username,
		PasswordHash: passwordHash,
		Session: db_dto.Session{
			RefreshToken: tokens.RefreshToken,
			ExpiresAt:    expiresAt,
		},
	}
	userID, err := u.userService.Create(ctx, dbDTO)
	if err != nil {
		return tokens, apperror.NewAppError(err, "failed to create user", "", "US-000007")
	}

	tokens.AccessToken, err = u.tokenManager.GenerateAccessToken(userID)
	if err != nil {
		return tokens, apperror.NewAppError(err, "unknown error", "failed to create jwt token", "TJ-000001")
	}
	return tokens, nil
}

func (u userUsecase) AuthUser(ctx context.Context, dto AuthUserDTO) (Tokens, error) {
	var tokens Tokens
	passwordHash := utils.HashPassword(dto.Password)
	userDTO := db_dto.AuthUserDTO{
		Email:        dto.Email,
		Username:     dto.Username,
		PasswordHash: passwordHash,
	}
	userID, err := u.userService.AuthUser(ctx, userDTO)
	if err != nil {
		return tokens, err
	}

	tokens.AccessToken, err = u.tokenManager.GenerateAccessToken(userID)
	if err != nil {
		return tokens, apperror.NewAppError(err, "unknown error", "failed to create jwt token", "TJ-000001")
	}
	tokens.RefreshToken, err = u.tokenManager.GenerateRefreshToken()
	if err != nil {
		return tokens, apperror.NewAppError(err, "unknown error", "failed to create refresh jwt token", "TJ-000001")
	}
	return tokens, nil
}

func (u userUsecase) GetUser(ctx context.Context, dto GetUserDTO) (user entity.User, err error) {
	userID, err := u.tokenManager.TokenUser(dto.Token)
	if err != nil {
		return user, apperror.NewAppError(err, "token error", err.Error(), "TJ-000002")
	}
	getDTO := db_dto.GetUserDTO{
		ID: userID,
	}
	dto.Token = userID
	user, err = u.userService.GetUser(ctx, getDTO)
	if err != nil {
		return user, apperror.NewAppError(err, "get user error", "", "US-000009")
	}
	return user, nil
}

func (u userUsecase) RefreshToken(ctx context.Context, dto RefreshTokenDTO) (Tokens, error) {
	var tokens Tokens
	ok, err := u.tokenManager.ValidateToken(dto.Token)
	if err != nil {
		return tokens, apperror.NewAppError(err, "refresh token error", "", "US-000005")
	}
	if !ok {
		return tokens, apperror.NewAppError(nil, "token is expired", "", "TJ-000005")
	}
	userDTO := db_dto.GetUserByRefreshTokenDTO{
		Token: dto.Token,
	}
	userID, err := u.userService.GetUserByRefreshToken(ctx, userDTO)
	if err != nil {
		return tokens, apperror.NewAppError(err, "refresh token error", "", "US-000005")
	}
	token, err := u.tokenManager.GenerateAccessToken(userID)
	if err != nil {
		return tokens, apperror.NewAppError(err, "unknown error", "failed to create jwt token", "TJ-000001")
	}
	tokens.AccessToken = token
	tokens.RefreshToken = dto.Token
	return tokens, nil
}

func NewUserUsecase(service Service, tokenManager TokenManager) *userUsecase {
	return &userUsecase{
		userService:  service,
		tokenManager: tokenManager,
	}
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
