package usecase

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrEmptyQuery     = fmt.Errorf("%w: query is required", ErrInvalidInput)
	ErrNegativeOffset = fmt.Errorf("%w: offset must not be negative", ErrInvalidInput)
	ErrNegativeLimit  = fmt.Errorf("%w: limit must not be negative", ErrInvalidInput)
)
