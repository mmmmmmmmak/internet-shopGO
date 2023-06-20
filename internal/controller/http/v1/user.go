package v1

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"main/internal/apperror"
	"main/internal/controller/http/dto"
	"main/internal/domain/entity"
	"main/internal/domain/usecase/user"
	"main/pkg/logging"
	"net/http"
	"strings"
)

const (
	usersURL = "/users"
	userURL  = "/users/:uuid"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, dto user_usecase.CreateUserDTO) (user_usecase.Tokens, error)
	AuthUser(ctx context.Context, dto user_usecase.AuthUserDTO) (user_usecase.Tokens, error)
	GetUser(ctx context.Context, dto user_usecase.GetUserDTO) (entity.User, error)
	RefreshToken(ctx context.Context, dto user_usecase.RefreshTokenDTO) (user_usecase.Tokens, error)
}

type userHandler struct {
	userUsecase  UserUsecase
	tokenManager TokenManager
	logger       *logging.Logger
}

func NewUserHandler(userUsecase UserUsecase, tokenManager TokenManager, logger *logging.Logger) *userHandler {
	return &userHandler{
		userUsecase:  userUsecase,
		tokenManager: tokenManager,
		logger:       logger,
	}
}

func (h *userHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/login/auth", apperror.Middleware(h.Auth))
	router.HandlerFunc(http.MethodPost, "/login/registration", apperror.Middleware(h.Registration))
	router.HandlerFunc(http.MethodGet, "/getUser", apperror.Middleware(h.GetUser))
	router.HandlerFunc(http.MethodPost, "/refreshToken", apperror.Middleware(h.RefreshToken))
}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) error {
	var d dto.GetUserDTO
	if val := r.Header.Get("Authorization"); val != "" {
		d.Token = strings.TrimPrefix(val, "Bearer ")
	} else {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return nil
	}

	usecaseDTO := user_usecase.GetUserDTO{
		Token: d.Token,
	}
	user, err := h.userUsecase.GetUser(r.Context(), usecaseDTO)
	if err != nil {
		// JSON RPC: TRANSPORT: 200, error: {msg, ..., dev_msg}
		return err
	}
	w.WriteHeader(http.StatusOK)
	jsonResponse, err := json.Marshal(user)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	return nil
}

func (h *userHandler) RefreshToken(w http.ResponseWriter, r *http.Request) error {
	var d dto.RefreshTokenDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		return err
	}

	usecaseDTO := user_usecase.RefreshTokenDTO{
		Token: d.Token,
	}
	tokens, err := h.userUsecase.RefreshToken(r.Context(), usecaseDTO)
	if err != nil {
		// JSON RPC: TRANSPORT: 200, error: {msg, ..., dev_msg}
		return err
	}
	w.WriteHeader(http.StatusOK)
	jsonResponse, err := json.Marshal(tokens)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	return nil
}
