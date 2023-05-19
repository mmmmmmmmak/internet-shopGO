package v1

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"main/internal/controller/http/dto"
	user_usecase "main/internal/domain/usecase/user"
	"main/pkg/apperror"
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
	w.Header().Set("Authorization", "Bearer "+user)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(user))

	return nil
}

func (h *userHandler) Auth(w http.ResponseWriter, r *http.Request) error {
	var d dto.Auth
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		return err
	}

	if d.Email == "" || !govalidator.IsEmail(d.Email) {
		if d.Username == "" {
			return apperror.NewAppError(nil, "username and email exists", "", "US-000008")
		}
		usecaseDTO := user_usecase.AuthByUsername{
			Username: d.Username,
			Password: d.Password,
		}
		user, err := h.userUsecase.AuthByUsername(r.Context(), usecaseDTO)
		if err != nil {
			// JSON RPC: TRANSPORT: 200, error: {msg, ..., dev_msg}
			return err
		}
		w.Header().Set("Authorization", "Bearer "+user)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(user))

		return nil
	}

	usecaseDTO := user_usecase.AuthByEmail{
		Email:    d.Email,
		Password: d.Password,
	}
	user, err := h.userUsecase.AuthByEmail(r.Context(), usecaseDTO)
	if err != nil {
		// JSON RPC: TRANSPORT: 200, error: {msg, ..., dev_msg}
		return err
	}
	w.Header().Set("Authorization", "Bearer "+user)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(user))

	return nil
}
