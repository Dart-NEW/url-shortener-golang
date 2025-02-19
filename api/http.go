package api

import (
	"encoding/json"
	"net/http"

	"url-shortener-golang/storage"
)

type ShortenRequest struct {
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

	var request ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
	defer r.Body.Close()

	shortURL := h.storage.Post(request.URL)

	response := map[string]string{"short_url": shortURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Resolve извлекает сокращённую ссылку из запроса и возвращает JSON оригинальную ссылку
func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Query().Get("short_url")

	if shortURL == "" {
		http.Error(w, "Missing or invalid 'short_url' parameter", http.StatusBadRequest)
		return
	}

	originalURL, exists := h.storage.Get(shortURL)
	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	response := map[string]string{"original_url": originalURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterRoutes регистрирует обработчики
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/post", h.Shorten)
	mux.HandleFunc("/get", h.Resolve)
}
