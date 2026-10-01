package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/hansenhalim/elelibrary/be/book/repository/googlebooks"
	"github.com/hansenhalim/elelibrary/be/book/usecase"
	"github.com/hansenhalim/elelibrary/be/handler/nethttp"
)

func main() {
	apiKey := os.Getenv("GOOGLE_BOOKS_API_KEY")
	if apiKey == "" {
		slog.Error("GOOGLE_BOOKS_API_KEY is not set")
		os.Exit(1)
	}

	bookSearcher := googlebooks.NewBookSearcher(
		&http.Client{Timeout: 10 * time.Second},
		googlebooks.BaseURL,
		apiKey,
	)
	h := nethttp.New(usecase.NewListBooks(bookSearcher))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/books", h.ListBooks)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("starting server", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
