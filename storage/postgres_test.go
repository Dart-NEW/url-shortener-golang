package storage_test

import (
	"testing"
	"url-shortener-golang/shortener"
	"url-shortener-golang/storage"

	"github.com/stretchr/testify/assert"
)

func TestPostgresStorage(t *testing.T) {
	connStr := "user=dart password=123 dbname=shortener sslmode=disable"

	st, err := storage.NewPostgresStorage(connStr)
	assert.NoError(t, err)
	defer st.Close()

	t.Run("Basic operations", func(t *testing.T) {
		url := "https://test-postgres.com"

		short, err := st.Post(url)
		assert.NoError(t, err)
		assert.Equal(t, short, shortener.GenerateShortURL(url))

		original, err := st.Get(short)
		assert.NoError(t, err)
		assert.Equal(t, url, original)

		short2, err := st.Post(url)
		assert.NoError(t, err)
		assert.Equal(t, short, short2)

		_, err = st.Get("nothing")
		assert.ErrorIs(t, err, storage.ErrNotFound)
	})
}
