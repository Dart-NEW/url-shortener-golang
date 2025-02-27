package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"url-shortener-golang/storage"
)

type ShortenRequestHttp struct {
	URL string `json:"url"`
}

type Handler struct {
	storage storage.Storage
}

// NewHandler создаёт новый экземпляр HTTP-обработчика
func NewHandler(storage storage.Storage) *Handler {
	return &Handler{storage: storage}
}

// Shorten обрабатывает JSON с оригинальной ссылкой и вовзращает JSON сокращённую ссылку
func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var request ShortenRequestHttp
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	shortURL, err := h.storage.Post(request.URL)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"short_url": shortURL}
	json.NewEncoder(w).Encode(response)
}

// Resolve извлекает сокращённую ссылку из запроса и возвращает JSON оригинальную ссылку
func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Query().Get("short_url")

	if shortURL == "" {
		http.Error(w, "Missing or invalid 'short_url' parameter", http.StatusBadRequest)
		return
	}

	originalURL, err := h.storage.Get(shortURL)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"original_url": originalURL}
	json.NewEncoder(w).Encode(response)
}

// RegisterRoutes регистрирует обработчики
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/post", h.Shorten)
	mux.HandleFunc("/get", h.Resolve)
}
