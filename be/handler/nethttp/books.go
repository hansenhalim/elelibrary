package nethttp

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/hansenhalim/elelibrary/be/book/usecase"
)

type listBooksResponse struct {
	Books []bookResponse `json:"books"`
}

type bookResponse struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Authors      []string `json:"authors"`
	ThumbnailURL string   `json:"thumbnailUrl"`
	Rating       float64  `json:"rating"`
	IsFavorite   bool     `json:"isFavorite"`
}

func (h *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	offset, err := intParam(query, "offset")
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "offset must be an integer")
		return
	}
	limit, err := intParam(query, "limit")
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "limit must be an integer")
		return
	}

	books, err := h.listBooks.Execute(r.Context(), usecase.ListBooksInput{
		Query:  query.Get("q"),
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}

	resp := listBooksResponse{Books: make([]bookResponse, 0, len(books))}
	for _, b := range books {
		resp.Books = append(resp.Books, toBookResponse(b))
	}
	writeJSON(w, http.StatusOK, resp)
}

func intParam(query url.Values, key string) (int, error) {
	v := query.Get(key)
	if v == "" {
		return 0, nil
	}
	return strconv.Atoi(v)
}

func toBookResponse(b usecase.ListedBook) bookResponse {
	authors := b.Book.Authors
	if authors == nil {
		authors = []string{}
	}
	return bookResponse{
		ID:           b.Book.ID,
		Title:        b.Book.Title,
		Authors:      authors,
		ThumbnailURL: b.Book.ThumbnailURL,
		Rating:       b.Book.Rating,
		IsFavorite:   b.IsFavorite,
	}
}
