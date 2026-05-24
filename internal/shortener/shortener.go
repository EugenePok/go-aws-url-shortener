package shortener

import (
	"crypto/sha256"

	"github.com/eugenepok/go-aws-url-shortener/internal/base62"
)

func canonicalize(url string) string {
	return url
}

func EncodeToShortURL(url string) string {
	canonicalizedUrl := canonicalize(url)
	hash := sha256.Sum256([]byte(canonicalizedUrl))
	return base62.StdEncoding.EncodeToString(hash[:])[:8]
}
