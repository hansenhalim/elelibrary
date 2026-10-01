package usecase

import (
	"context"
	"fmt"

	"github.com/hansenhalim/elelibrary/be/entity"
)

type FavoriteAdder interface {
	AddFavorite(ctx context.Context, bookID string) error
}

type AddFavorite struct {
	adder FavoriteAdder
}

func NewAddFavorite(adder FavoriteAdder) *AddFavorite {
	return &AddFavorite{adder: adder}
}

func (uc *AddFavorite) Execute(ctx context.Context, bookID string) error {
	if !entity.IsValidBookID(bookID) {
		return ErrInvalidBookID
	}
	if err := uc.adder.AddFavorite(ctx, bookID); err != nil {
		return fmt.Errorf("add favorite: %w", err)
	}
	return nil
}
