package entity

import "regexp"

var bookIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

type Book struct {
	ID           string
	Title        string
	Authors      []string
	ThumbnailURL string
	Rating       float64
}

func IsValidBookID(id string) bool {
	return bookIDPattern.MatchString(id)
}
