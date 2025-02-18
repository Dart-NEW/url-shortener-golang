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
func (s *MemoryStorage) Post(originalURL string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	for short, original := range s.urls {
		if original == originalURL {
			return short
		}
	}

	shortURL := shortener.GenerateShortURL(originalURL)
	s.urls[shortURL] = originalURL
	return shortURL
}

// Get возвращает оригиналный URL по сокращённому
func (s *MemoryStorage) Get(shortURL string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	originalUrl, exists := s.urls[shortURL]
	return originalUrl, exists
}

// Временная generateShortURL для создания короткого URL
// func generateShortURL(originalURL string) string {
// 	return "short_" + originalURL[:5]
// }
