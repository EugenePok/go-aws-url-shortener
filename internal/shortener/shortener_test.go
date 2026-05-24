package shortener

import "testing"

func TestEncodeSuccess(t *testing.T) {
	t.Run("Test Google", func(t *testing.T) {
		got := EncodeToShortURL("https://www.google.com")
		want := "eswp9Xga"
		if got != want {
			t.Errorf(" got %v, want %v", got, want)
		}
	})
	t.Run("Test Instagram", func(t *testing.T) {
		got := EncodeToShortURL("https://www.instagram.com")
		want := "aVCrgKYt"
		if got != want {
			t.Errorf(" got %v, want %v", got, want)
		}
	})
}
