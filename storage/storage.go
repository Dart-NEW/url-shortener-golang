package storage

import "errors"

var (
	ErrNotFound = errors.New("URL not found")
)

type Storage interface {
	Post(originalURL string) (string, error)
	Get(shortURL string) (string, error)
	Close() error
}
