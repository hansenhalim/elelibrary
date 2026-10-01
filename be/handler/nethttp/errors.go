package nethttp

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/hansenhalim/elelibrary/be/book/usecase"
)

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, usecase.ErrInvalidInput):
		writeMessage(w, http.StatusBadRequest, err.Error())
	default:
		slog.ErrorContext(r.Context(), "request failed", "method", r.Method, "path", r.URL.Path, "error", err)
		writeMessage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}
