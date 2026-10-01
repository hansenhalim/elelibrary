package usecase

import (
	"context"
	"fmt"
)

type FavoriteRemover interface {
	RemoveFavorite(ctx context.Context, bookID string) error
}

type RemoveFavorite struct {
	remover FavoriteRemover
}

func NewRemoveFavorite(remover FavoriteRemover) *RemoveFavorite {
	return &RemoveFavorite{remover: remover}
}

func (uc *RemoveFavorite) Execute(ctx context.Context, bookID string) error {
	if err := validateBookID(bookID); err != nil {
		return err
	}
	if err := uc.remover.RemoveFavorite(ctx, bookID); err != nil {
		return fmt.Errorf("remove favorite: %w", err)
	}
	return nil
}
