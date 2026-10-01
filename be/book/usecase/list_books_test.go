package usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hansenhalim/elelibrary/be/book/usecase"
	"github.com/hansenhalim/elelibrary/be/entity"
)

func TestListBooks_Execute(t *testing.T) {
	book := entity.Book{ID: "zyTCAlFPjgYC", Title: "The Google Story"}

	tests := []struct {
		name       string
		in         usecase.ListBooksInput
		wantQuery  string
		wantOffset int
		wantLimit  int
	}{
		{
			name:       "passes input to searcher",
			in:         usecase.ListBooksInput{Query: "go", Offset: 20, Limit: 15},
			wantQuery:  "go",
			wantOffset: 20,
			wantLimit:  15,
		},
		{
			name:      "trims query",
			in:        usecase.ListBooksInput{Query: "  go  ", Limit: 10},
			wantQuery: "go",
			wantLimit: 10,
		},
		{
			name:      "defaults zero limit to 10",
			in:        usecase.ListBooksInput{Query: "go"},
			wantQuery: "go",
			wantLimit: 10,
		},
		{
			name:      "caps limit at 20",
			in:        usecase.ListBooksInput{Query: "go", Limit: 100},
			wantQuery: "go",
			wantLimit: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searcher := usecase.NewMockBookSearcher(t)
			searcher.EXPECT().
				SearchBooks(mock.Anything, tt.wantQuery, tt.wantOffset, tt.wantLimit).
				Return([]entity.Book{book}, nil)
			finder := usecase.NewMockFavoriteFinder(t)
			finder.EXPECT().
				FindFavoriteIDs(mock.Anything, []string{book.ID}).
				Return(nil, nil)

			got, err := usecase.NewListBooks(searcher, finder).Execute(t.Context(), tt.in)

			require.NoError(t, err)
			assert.Equal(t, []usecase.ListedBook{{Book: book}}, got)
		})
	}
}

func TestListBooks_Execute_InvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		in      usecase.ListBooksInput
		wantErr error
	}{
		{
			name:    "empty query",
			in:      usecase.ListBooksInput{Query: ""},
			wantErr: usecase.ErrEmptyQuery,
		},
		{
			name:    "blank query",
			in:      usecase.ListBooksInput{Query: "   "},
			wantErr: usecase.ErrEmptyQuery,
		},
		{
			name:    "negative offset",
			in:      usecase.ListBooksInput{Query: "go", Offset: -1},
			wantErr: usecase.ErrNegativeOffset,
		},
		{
			name:    "negative limit",
			in:      usecase.ListBooksInput{Query: "go", Limit: -1},
			wantErr: usecase.ErrNegativeLimit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searcher := usecase.NewMockBookSearcher(t)
			finder := usecase.NewMockFavoriteFinder(t)

			got, err := usecase.NewListBooks(searcher, finder).Execute(t.Context(), tt.in)

			require.ErrorIs(t, err, tt.wantErr)
			assert.ErrorIs(t, err, usecase.ErrInvalidInput)
			assert.Nil(t, got)
		})
	}
}

func TestListBooks_Execute_SearcherError(t *testing.T) {
	errUpstream := errors.New("google books unavailable")
	searcher := usecase.NewMockBookSearcher(t)
	searcher.EXPECT().
		SearchBooks(mock.Anything, "go", 0, 10).
		Return(nil, errUpstream)
	finder := usecase.NewMockFavoriteFinder(t)

	got, err := usecase.NewListBooks(searcher, finder).Execute(t.Context(), usecase.ListBooksInput{Query: "go"})

	require.ErrorIs(t, err, errUpstream)
	assert.NotErrorIs(t, err, usecase.ErrInvalidInput)
	assert.Nil(t, got)
}

func TestListBooks_Execute_MarksFavorites(t *testing.T) {
	books := []entity.Book{{ID: "a"}, {ID: "b"}, {ID: "c"}}
	searcher := usecase.NewMockBookSearcher(t)
	searcher.EXPECT().
		SearchBooks(mock.Anything, "go", 0, 10).
		Return(books, nil)
	finder := usecase.NewMockFavoriteFinder(t)
	finder.EXPECT().
		FindFavoriteIDs(mock.Anything, []string{"a", "b", "c"}).
		Return([]string{"c", "a"}, nil)

	got, err := usecase.NewListBooks(searcher, finder).Execute(t.Context(), usecase.ListBooksInput{Query: "go"})

	require.NoError(t, err)
	assert.Equal(t, []usecase.ListedBook{
		{Book: books[0], IsFavorite: true},
		{Book: books[1], IsFavorite: false},
		{Book: books[2], IsFavorite: true},
	}, got)
}

func TestListBooks_Execute_NoBooksSkipsFavorites(t *testing.T) {
	searcher := usecase.NewMockBookSearcher(t)
	searcher.EXPECT().
		SearchBooks(mock.Anything, "go", 0, 10).
		Return(nil, nil)
	finder := usecase.NewMockFavoriteFinder(t)

	got, err := usecase.NewListBooks(searcher, finder).Execute(t.Context(), usecase.ListBooksInput{Query: "go"})

	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestListBooks_Execute_FinderError(t *testing.T) {
	errDB := errors.New("connection refused")
	searcher := usecase.NewMockBookSearcher(t)
	searcher.EXPECT().
		SearchBooks(mock.Anything, "go", 0, 10).
		Return([]entity.Book{{ID: "a"}}, nil)
	finder := usecase.NewMockFavoriteFinder(t)
	finder.EXPECT().
		FindFavoriteIDs(mock.Anything, []string{"a"}).
		Return(nil, errDB)

	got, err := usecase.NewListBooks(searcher, finder).Execute(t.Context(), usecase.ListBooksInput{Query: "go"})

	require.ErrorIs(t, err, errDB)
	assert.NotErrorIs(t, err, usecase.ErrInvalidInput)
	assert.Nil(t, got)
}
