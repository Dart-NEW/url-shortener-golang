package storage

import (
	"sync"
	"url-shortener-golang/shortener"
)

type MemoryStorage struct {
	mu   sync.RWMutex
	urls map[string]string
}

// NewMemoryStorage для создания нового экземпляра MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		urls: make(map[string]string),
	}
}

// Post для сохранения оригинального URL и возврата сокращённого
func (s *MemoryStorage) Post(originalURL string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for short, original := range s.urls {
		if original == originalURL {
			return short, nil
		}
	}

	shortURL := shortener.GenerateShortURL(originalURL)
	s.urls[shortURL] = originalURL
	return shortURL, nil
}

// Get возвращает оригиналный URL по сокращённому
func (s *MemoryStorage) Get(shortURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	originalUrl, exists := s.urls[shortURL]
	if !exists {
		return "", ErrNotFound
	}
	return originalUrl, nil
}

// Close закрывает соединение с хранилищем (пустой метод для in-memory).
func (s *MemoryStorage) Close() error {
	return nil
}
