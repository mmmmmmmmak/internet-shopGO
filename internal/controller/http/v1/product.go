package v1

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"main/internal/apperror"
	"main/internal/controller/http/dto"
	"main/internal/domain/entity"
	product_usecase "main/internal/domain/usecase/product"
	"strings"

	//user_usecase "main/internal/domain/usecase/user"
	"main/pkg/logging"
	"net/http"
)

type ProductUsecase interface {
	GetAll(ctx context.Context) ([]entity.Product, error)
	AddProduct(ctx context.Context, dto product_usecase.CreateProductDTO) (string, error)
}

type TokenManager interface {
	CheckIsSeller(next http.HandlerFunc) http.HandlerFunc
}

type productHandler struct {
	productUsecase ProductUsecase
	tokenManager   TokenManager
	logger         *logging.Logger
}

func NewProductHandler(productUsecase ProductUsecase, tokenManager TokenManager, logger *logging.Logger) *productHandler {
	return &productHandler{
		productUsecase: productUsecase,
		tokenManager:   tokenManager,
		logger:         logger,
	}
}

func (h *productHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/product", apperror.Middleware(h.GetAll))
	router.HandlerFunc(http.MethodPost, "/product", h.tokenManager.CheckIsSeller(apperror.Middleware(h.AddProduct)))
}

func (h *productHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	products, err := h.productUsecase.GetAll(r.Context())
	if err != nil {
		// JSON RPC: TRANSPORT: 200, error: {msg, ..., dev_msg}
		return err
	}
	w.WriteHeader(http.StatusOK)
	jsonResponse, err := json.Marshal(products)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
	return nil
}

func (h *productHandler) AddProduct(w http.ResponseWriter, r *http.Request) error {
	var d dto.CreateProductDTO
	if val := r.Header.Get("Authorization"); val != "" {
		d.Token = strings.TrimPrefix(val, "Bearer ")
	} else {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return nil
	}
	response := make(map[string]string)
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		return err
	}
	usecaseDTO := product_usecase.CreateProductDTO{
		Name:        d.Name,
		Description: d.Description,
		Price:       d.Price,
		Token:       d.Token,
	}
	idProduct, err := h.productUsecase.AddProduct(r.Context(), usecaseDTO)
	if err != nil {
		return err
	}
	response["id"] = idProduct
	w.WriteHeader(http.StatusOK)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
	return nil
}
