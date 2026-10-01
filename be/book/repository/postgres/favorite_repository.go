package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hansenhalim/elelibrary/be/book/usecase"
)

var (
	_ usecase.FavoriteFinder  = (*FavoriteRepository)(nil)
	_ usecase.FavoriteAdder   = (*FavoriteRepository)(nil)
	_ usecase.FavoriteRemover = (*FavoriteRepository)(nil)
)

type FavoriteRepository struct {
	pool *pgxpool.Pool
}

func NewFavoriteRepository(pool *pgxpool.Pool) *FavoriteRepository {
	return &FavoriteRepository{pool: pool}
}

func (r *FavoriteRepository) FindFavoriteIDs(ctx context.Context, bookIDs []string) ([]string, error) {
	rows, err := r.pool.Query(ctx, `SELECT book_id FROM favorites WHERE book_id = ANY($1)`, bookIDs)
	if err != nil {
		return nil, fmt.Errorf("postgres: find favorite ids: %w", err)
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, fmt.Errorf("postgres: find favorite ids: %w", err)
	}
	return ids, nil
}

func (r *FavoriteRepository) AddFavorite(ctx context.Context, bookID string) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO favorites (book_id) VALUES ($1) ON CONFLICT (book_id) DO NOTHING`, bookID)
	if err != nil {
		return fmt.Errorf("postgres: add favorite: %w", err)
	}
	return nil
}

func (r *FavoriteRepository) RemoveFavorite(ctx context.Context, bookID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM favorites WHERE book_id = $1`, bookID)
	if err != nil {
		return fmt.Errorf("postgres: remove favorite: %w", err)
	}
	return nil
}
