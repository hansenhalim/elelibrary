package usecase_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hansenhalim/elelibrary/be/book/usecase"
)

func TestRemoveFavorite_Execute(t *testing.T) {
	remover := usecase.NewMockFavoriteRemover(t)
	remover.EXPECT().RemoveFavorite(mock.Anything, "zyTCAlFPjgYC").Return(nil)

	err := usecase.NewRemoveFavorite(remover).Execute(t.Context(), "zyTCAlFPjgYC")

	require.NoError(t, err)
}

func TestRemoveFavorite_Execute_InvalidBookID(t *testing.T) {
	tests := []struct {
		name   string
		bookID string
	}{
		{name: "empty", bookID: ""},
		{name: "disallowed characters", bookID: "abc/def"},
		{name: "too long", bookID: strings.Repeat("a", 65)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remover := usecase.NewMockFavoriteRemover(t)

			err := usecase.NewRemoveFavorite(remover).Execute(t.Context(), tt.bookID)

			require.ErrorIs(t, err, usecase.ErrInvalidBookID)
			assert.ErrorIs(t, err, usecase.ErrInvalidInput)
		})
	}
}

func TestRemoveFavorite_Execute_RemoverError(t *testing.T) {
	errDB := errors.New("connection refused")
	remover := usecase.NewMockFavoriteRemover(t)
	remover.EXPECT().RemoveFavorite(mock.Anything, "zyTCAlFPjgYC").Return(errDB)

	err := usecase.NewRemoveFavorite(remover).Execute(t.Context(), "zyTCAlFPjgYC")

	require.ErrorIs(t, err, errDB)
	assert.NotErrorIs(t, err, usecase.ErrInvalidInput)
}
