package apierror

import (
	"errors"
	"net/http"
)

type MyHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func Middleware(next MyHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var apiErr *ApiError
		err := next(w, r)
		if err != nil {
			if errors.As(err, &apiErr) {
				if errors.Is(err, ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					w.Write(ErrNotFound.Marshal())
					return
				}
				err := err.(*ApiError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(err.Marshal())
				return
			}
			w.WriteHeader(http.StatusTeapot)
			w.Write(systemErr(err.Error()).Marshal())
		}
	}
}
