package v1

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"main/internal/apperror"
	"main/internal/controller/http/dto"
	user_usecase "main/internal/domain/usecase/user"
	"net/http"
)

func (h *userHandler) Registration(w http.ResponseWriter, r *http.Request) error {
	var d dto.CreateUserDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		return err
	}

	// validate

	// MAPPING dto.CreateBookDTO --> book_usecase.CreateBookDTO
	usecaseDTO := user_usecase.CreateUserDTO{
		Email:    d.Email,
		Username: d.Username,
		Password: d.Password,
	}
	user, err := h.userUsecase.CreateUser(r.Context(), usecaseDTO)
	if err != nil {
		// JSON RPC: TRANSPORT: 200, error: {msg, ..., dev_msg}
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(user))

	return nil
}

func (h *userHandler) Auth(w http.ResponseWriter, r *http.Request) error {
	response := make(map[string]string)
	var d dto.Auth
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		return err
	}

	if (d.Email == "" && d.Username == "") || (!govalidator.IsEmail(d.Email) && d.Username == "") {
		return apperror.NewAppError(nil, "incorrect data entered", "", "US-000004")
	}
	usecaseDTO := user_usecase.AuthUser{
		Email:    d.Email,
		Username: d.Username,
		Password: d.Password,
	}
	user, err := h.userUsecase.AuthUser(r.Context(), usecaseDTO)
	if err != nil {
		// JSON RPC: TRANSPORT: 200, error: {msg, ..., dev_msg}
		return err
	}
	w.WriteHeader(http.StatusOK)
	response["jwt"] = user
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	return nil
}
