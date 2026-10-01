package usecase

import "regexp"

var bookIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

func validateBookID(id string) error {
	if !bookIDPattern.MatchString(id) {
		return ErrInvalidBookID
	}
	return nil
}
