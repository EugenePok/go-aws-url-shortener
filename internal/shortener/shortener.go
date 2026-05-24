package shortener

import (
	"crypto/sha256"
	"time"

	"github.com/eugenepok/go-aws-url-shortener/internal/base62"
	"github.com/eugenepok/go-aws-url-shortener/pkg/models"
)

func canonicalize(url string) string {
	return url
}

func encodeToShortURL(url string) string {
	canonicalizedUrl := canonicalize(url)
	hash := sha256.Sum256([]byte(canonicalizedUrl))
	return base62.StdEncoding.EncodeToString(hash[:])[:8]
}

func Encode(url string) *models.UrlData {
	return &models.UrlData{
		FullURL:   url,
		ShortURL:  encodeToShortURL(url),
		CreatedAt: time.Now(),
	}
}
