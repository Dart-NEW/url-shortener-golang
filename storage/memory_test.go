package storage_test

import (
	"testing"
	"url-shortener-golang/storage"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage(t *testing.T) {
	s := storage.NewMemoryStorage()

	t.Run("Post and get URL", func(t *testing.T) {
		url := "https://example.com"
		short, err := s.Post(url)
		assert.NoError(t, err)

		original, err := s.Get(short)
		assert.NoError(t, err)
		assert.Equal(t, original, url)
	})

	t.Run("Same URL returns same short", func(t *testing.T) {
		url := "https://same_urls.com"
		short1, err := s.Post(url)
		assert.NoError(t, err)
		short2, err := s.Post(url)
		assert.NoError(t, err)
		assert.Equal(t, short1, short2)
	})

	t.Run("Get non-exsistent URL", func(t *testing.T) {
		url := "nothing"
		_, err := s.Get(url)
		assert.ErrorIs(t, err, storage.ErrNotFound)
	})
}
