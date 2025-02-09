package tests

import (
	"math/rand"
	"strings"
	"testing"

	urlshortener "example.com/m/internal/url_shortener"
	"github.com/stretchr/testify/assert"
)

func TestUrlShortener(t *testing.T) {
	t.Run("returns string [A-Za-z0-9_]", func(t *testing.T) {
		type testCase struct {
			str      string
			expected bool
		}

		testCases := []testCase{
			{
				str:      urlshortener.Shorten(12412),
				expected: true,
			},
			{
				str:      urlshortener.Shorten(992132198),
				expected: true,
			},
			{
				str:      urlshortener.Shorten(0),
				expected: true,
			},
		}

		for _, tc := range testCases {
			assert.Equal(t, tc.expected, validateString(tc.str))
		}
	})

	t.Run("is idempotent", func(t *testing.T) {
		num := rand.Uint32()
		str := urlshortener.Shorten(num)

		for i := 0; i < 10000; i++ {
			assert.Equal(t, str, urlshortener.Shorten(num))
		}
	})

	t.Run("colision test", func(t *testing.T) {
		seen := make(map[string]bool)
		for i := 0; i < 10000; i++ {
			s := urlshortener.Shorten(rand.Uint32())
			if seen[s] {
				t.Fatalf("Collision detected for %d: %s", i, s)
			}
			seen[s] = true
		}
	})
}

func validateString(str string) bool {
	if len(str) != 10 {
		return false
	}

	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

	for i := 0; i < 10; i++ {
		if !strings.Contains(alphabet, string(str[i])) {
			return false
		}
	}

	return true
}
