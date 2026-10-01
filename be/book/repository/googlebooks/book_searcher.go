package googlebooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hansenhalim/elelibrary/be/book/usecase"
	"github.com/hansenhalim/elelibrary/be/entity"
)

const BaseURL = "https://www.googleapis.com/books/v1"

const maxResultWindow = 2000

var _ usecase.BookSearcher = (*BookSearcher)(nil)

type BookSearcher struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

func NewBookSearcher(client *http.Client, baseURL, apiKey string) *BookSearcher {
	return &BookSearcher{client: client, baseURL: baseURL, apiKey: apiKey}
}

func (s *BookSearcher) SearchBooks(ctx context.Context, query string, offset, limit int) ([]entity.Book, error) {
	if offset >= maxResultWindow {
		return nil, nil
	}
	limit = min(limit, maxResultWindow-offset)

	params := url.Values{
		"q":          {query},
		"startIndex": {strconv.Itoa(offset)},
		"maxResults": {strconv.Itoa(limit)},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/volumes?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("google books: build request: %w", err)
	}
	req.Header.Set("X-Goog-Api-Key", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("google books: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, statusError(resp)
	}

	var body volumesResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("google books: decode response: %w", err)
	}

	books := make([]entity.Book, 0, len(body.Items))
	for _, v := range body.Items {
		books = append(books, v.toEntity())
	}
	return books, nil
}

func statusError(resp *http.Response) error {
	var body errorResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil || body.Error.Message == "" {
		return fmt.Errorf("google books: unexpected status %s", resp.Status)
	}
	return fmt.Errorf("google books: unexpected status %s: %s", resp.Status, body.Error.Message)
}

func secureURL(u string) string {
	if rest, ok := strings.CutPrefix(u, "http://"); ok {
		return "https://" + rest
	}
	return u
}
