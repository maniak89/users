package storage

import (
	"context"
	"errors"

	"users/internal/models/storage"
)

var (
	ErrNoRows       = errors.New("no rows")
	ErrInvalidState = errors.New("invalid state")
)

type Storage interface {
	Find(ctx context.Context, login string) (*storage.User, error)
	UpdateLastLogin(ctx context.Context, login string) error
}
