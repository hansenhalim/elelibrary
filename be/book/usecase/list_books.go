package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/hansenhalim/elelibrary/be/entity"
)

const (
	defaultLimit = 10
	maxLimit     = 20
)

type BookSearcher interface {
	SearchBooks(ctx context.Context, query string, offset, limit int) ([]entity.Book, error)
}

type ListBooksInput struct {
	Query  string
	Offset int
	Limit  int
}

type ListBooks struct {
	searcher BookSearcher
}

func NewListBooks(searcher BookSearcher) *ListBooks {
	return &ListBooks{searcher: searcher}
}

func (uc *ListBooks) Execute(ctx context.Context, in ListBooksInput) ([]entity.Book, error) {
	query := strings.TrimSpace(in.Query)
	if query == "" {
		return nil, ErrEmptyQuery
	}
	if in.Offset < 0 {
		return nil, ErrNegativeOffset
	}
	if in.Limit < 0 {
		return nil, ErrNegativeLimit
	}

	limit := in.Limit
	if limit == 0 {
		limit = defaultLimit
	}
	limit = min(limit, maxLimit)

	books, err := uc.searcher.SearchBooks(ctx, query, in.Offset, limit)
	if err != nil {
		return nil, fmt.Errorf("search books: %w", err)
	}
	return books, nil
}
