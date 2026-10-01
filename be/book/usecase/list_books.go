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

type FavoriteFinder interface {
	FindFavoriteIDs(ctx context.Context, bookIDs []string) ([]string, error)
}

type ListBooksInput struct {
	Query  string
	Offset int
	Limit  int
}

type ListedBook struct {
	Book       entity.Book
	IsFavorite bool
}

type ListBooks struct {
	searcher BookSearcher
	finder   FavoriteFinder
}

func NewListBooks(searcher BookSearcher, finder FavoriteFinder) *ListBooks {
	return &ListBooks{searcher: searcher, finder: finder}
}

func (uc *ListBooks) Execute(ctx context.Context, in ListBooksInput) ([]ListedBook, error) {
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
	if len(books) == 0 {
		return nil, nil
	}

	ids := make([]string, 0, len(books))
	for _, b := range books {
		ids = append(ids, b.ID)
	}
	favoriteIDs, err := uc.finder.FindFavoriteIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("find favorite ids: %w", err)
	}
	favorites := make(map[string]bool, len(favoriteIDs))
	for _, id := range favoriteIDs {
		favorites[id] = true
	}

	listed := make([]ListedBook, 0, len(books))
	for _, b := range books {
		listed = append(listed, ListedBook{Book: b, IsFavorite: favorites[b.ID]})
	}
	return listed, nil
}
