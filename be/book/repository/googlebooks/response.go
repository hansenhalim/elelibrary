package googlebooks

import "github.com/hansenhalim/elelibrary/be/entity"

type volumesResponse struct {
	Items []volume `json:"items"`
}

type volume struct {
	ID         string     `json:"id"`
	VolumeInfo volumeInfo `json:"volumeInfo"`
}

type volumeInfo struct {
	Title         string     `json:"title"`
	Authors       []string   `json:"authors"`
	AverageRating float64    `json:"averageRating"`
	ImageLinks    imageLinks `json:"imageLinks"`
}

type imageLinks struct {
	Thumbnail string `json:"thumbnail"`
}

type errorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (v volume) toEntity() entity.Book {
	return entity.Book{
		ID:           v.ID,
		Title:        v.VolumeInfo.Title,
		Authors:      v.VolumeInfo.Authors,
		ThumbnailURL: secureURL(v.VolumeInfo.ImageLinks.Thumbnail),
		Rating:       v.VolumeInfo.AverageRating,
	}
}
