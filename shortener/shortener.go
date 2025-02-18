package shortener

import (
	"crypto/md5"
	"encoding/base64"
	"strings"
)

const (
	alphabet  = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"
	newLength = 10
)

var (
	base64Encoding = base64.NewEncoding(alphabet + "-").WithPadding(base64.NoPadding)
)

// GenerateShortURL создаёт короткую ссылку из 10 символов из оригинального URL
func GenerateShortURL(originalURL string) string {
	hash := md5.Sum([]byte(originalURL))

	encoded := base64Encoding.EncodeToString(hash[:])
	encoded = encoded[:newLength]

	shortURL := strings.Map(func(r rune) rune {
		if r == '-' {
			return '_'
		}
		return r
	}, encoded)

	return shortURL
}
