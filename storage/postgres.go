package storage

import (
	"database/sql"
	"errors"
	"url-shortener-golang/shortener"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage создаёт новое хранилище на основе PostgreSQL
func NewPostgresStorage(connectionString string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			short_url VARCHAR(10) PRIMARY KEY,
			original_url TEXT UNIQUE NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

// Post сохраняет оригинальный URL и возвращает сокращённый
func (s *PostgresStorage) Post(originalURL string) (string, error) {
	var existingShort string
	err := s.db.QueryRow("SELECT short_url FROM urls WHERE original_url = $1", originalURL).Scan(&existingShort)
	if err == nil {
		return existingShort, nil
	}

	shortURL := shortener.GenerateShortURL(originalURL)
	_, err = s.db.Exec("INSERT INTO urls (short_url, original_url) VALUES ($1, $2)", shortURL, originalURL)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

// Get возвращает оригинальный URL по сокращённому.
func (s *PostgresStorage) Get(shortURL string) (string, error) {
	var originalURL string
	err := s.db.QueryRow("SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}

	return originalURL, nil
}

// Close закрывает соединение с базой данных.
func (s *PostgresStorage) Close() error {
	return s.db.Close()
}
