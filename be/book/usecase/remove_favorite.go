package usecase

import (
	"context"
	"fmt"

	"github.com/hansenhalim/elelibrary/be/entity"
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
	if !entity.IsValidBookID(bookID) {
		return ErrInvalidBookID
	}
	if err := uc.remover.RemoveFavorite(ctx, bookID); err != nil {
		return fmt.Errorf("remove favorite: %w", err)
	}
	return nil
}
