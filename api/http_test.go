package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener-golang/api"
	"url-shortener-golang/storage"
)

func TestApiEndpoints(t *testing.T) {
	st := storage.NewMemoryStorage()
	h := api.NewHandler(st)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	t.Run("Valid shorten request", func(t *testing.T) {
		body := bytes.NewBufferString(`{"url":"https://valid.com"}`)
		req := httptest.NewRequest("POST", "/post", body)
		resp := httptest.NewRecorder()
		mux.ServeHTTP(resp, req)
		if resp.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.Code)
		}
		var result map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("Invalid shorten request", func(t *testing.T) {
		cases := []struct {
			name       string
			body       string
			statusCode int
		}{
			{"Empty body", "", http.StatusBadRequest},
			{"Invalid JSON", "{invalid}", http.StatusBadRequest},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				body := bytes.NewBufferString(tc.body)
				req := httptest.NewRequest("POST", "/post", body)
				resp := httptest.NewRecorder()

				mux.ServeHTTP(resp, req)

				if resp.Code != tc.statusCode {
					t.Errorf("Expected status %d, got %d", tc.statusCode, resp.Code)
				}
			})
		}
	})
	t.Run("Valid resolve request", func(t *testing.T) {
		short, _ := st.Post("https://test-resolve.com")

		req := httptest.NewRequest("GET", "/get?short_url="+short, nil)
		resp := httptest.NewRecorder()

		mux.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.Code)
		}

		var result map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}

		if result["original_url"] != "https://test-resolve.com" {
			t.Error("Invalid original URL returned")
		}
	})
	t.Run("Invalid resolve request", func(t *testing.T) {
		cases := []struct {
			name       string
			body       string
			statusCode int
		}{
			{"Missing short url", "", http.StatusBadRequest},
			{"Invalid short url", "url_short=https://test-resolve.com", http.StatusBadRequest},
			{"URL not found", "short_url=nothing", http.StatusNotFound},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest("POST", "/get?"+tc.body, nil)
				resp := httptest.NewRecorder()

				mux.ServeHTTP(resp, req)

				if resp.Code != tc.statusCode {
					t.Errorf("Expected status %d, got %d", tc.statusCode, resp.Code)
				}
			})
		}
	})
}
