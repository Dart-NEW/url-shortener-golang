package shortener_test

import (
	"testing"
	"url-shortener-golang/shortener"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Empty string", ""},
		{"Simple URL", "https://example.com"},
		{"Long URL", "https://super-duper-duper-mega-long-url-with-AIOLLDFRDLRAGSPIGdasfsaSHISLTIS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shortener.GenerateShortURL(tt.input)
			assert.Len(t, result, 10)
			assert.Regexp(t, `^[A-Za-z0-9_]`, result)
		})
	}
}
