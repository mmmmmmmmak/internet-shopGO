package apperror

import (
	"errors"
	"net/http"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var appError *AppError
		err := h(w, r)
		if err != nil {
			w.Header().Set("Content-Type", "application/type")
			if errors.As(err, &appError) {
				if errors.Is(err, ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					w.Write(ErrNotFound.Marshal())
					return
				}
				err = err.(*AppError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(appError.Marshal())
				return
			}
			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err).Marshal())
		}
	}
}
