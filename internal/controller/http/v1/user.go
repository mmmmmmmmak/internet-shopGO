package v1

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"main/internal/domain/usecase/user"
	"main/pkg/apperror"
	"main/pkg/logging"
	"net/http"
)

const (
	usersURL = "/users"
	userURL  = "/users/:uuid"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, dto user_usecase.CreateUserDTO) (string, error)
	AuthByEmail(ctx context.Context, dto user_usecase.AuthByEmail) (string, error)
	AuthByUsername(ctx context.Context, dto user_usecase.AuthByUsername) (string, error)
}

type userHandler struct {
	userUsecase UserUsecase
	logger      *logging.Logger
}

func NewUserHandler(userUsecase UserUsecase, logger *logging.Logger) *userHandler {
	return &userHandler{
		userUsecase: userUsecase,
		logger:      logger,
	}
}

func (h *userHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/login/auth", apperror.Middleware(h.Auth))
	router.HandlerFunc(http.MethodPost, "/login/registration", apperror.Middleware(h.Registration))
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.GetUserByUUID))
	router.HandlerFunc(http.MethodPut, userURL, apperror.Middleware(h.UpdateUser))
}

func (h *userHandler) GetUserByUUID(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(204)
	w.Write([]byte("this is update user"))
	return nil
}

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(204)
	w.Write([]byte("this is update user"))

	return nil
}
