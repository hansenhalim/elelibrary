package nethttp

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/hansenhalim/elelibrary/be/book/usecase"
)

type Handler struct {
	listBooks      *usecase.ListBooks
	addFavorite    *usecase.AddFavorite
	removeFavorite *usecase.RemoveFavorite
}

func New(listBooks *usecase.ListBooks, addFavorite *usecase.AddFavorite, removeFavorite *usecase.RemoveFavorite) *Handler {
	return &Handler{listBooks, addFavorite, removeFavorite}
}

type messageResponse struct {
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to write response", "error", err)
	}
}

func writeMessage(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, messageResponse{Message: message})
}
