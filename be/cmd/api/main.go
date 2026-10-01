package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hansenhalim/elelibrary/be/book/repository/googlebooks"
	"github.com/hansenhalim/elelibrary/be/book/repository/postgres"
	"github.com/hansenhalim/elelibrary/be/book/usecase"
	"github.com/hansenhalim/elelibrary/be/handler/nethttp"
)

func main() {
	apiKey := mustGetenv("GOOGLE_BOOKS_API_KEY")
	dbURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(mustGetenv("DB_USERNAME"), mustGetenv("DB_PASSWORD")),
		Host:   net.JoinHostPort(mustGetenv("DB_HOST"), mustGetenv("DB_PORT")),
		Path:   mustGetenv("DB_DATABASE"),
	}

	pool, err := pgxpool.New(context.Background(), dbURL.String())
	if err != nil {
		slog.Error("failed to create database pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	bookSearcher := googlebooks.NewBookSearcher(
		&http.Client{Timeout: 10 * time.Second},
		googlebooks.BaseURL,
		apiKey,
	)
	favorites := postgres.NewFavoriteRepository(pool)
	h := nethttp.New(
		usecase.NewListBooks(bookSearcher, favorites),
		usecase.NewAddFavorite(favorites),
		usecase.NewRemoveFavorite(favorites),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/books", h.ListBooks)
	mux.HandleFunc("PUT /api/favorites/{bookId}", h.AddFavorite)
	mux.HandleFunc("DELETE /api/favorites/{bookId}", h.RemoveFavorite)

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

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		slog.Error(key + " is not set")
		os.Exit(1)
	}
	return v
}
