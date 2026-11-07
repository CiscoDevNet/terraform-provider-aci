package acceptancetest

import (
	"crypto/rand"
	"math/big"
	"testing"
)

// PrefixOptions controls the characters and length of a fixture prefix.
type PrefixOptions struct {
	Characters string
	Length     int
}

// Prefix creates one random prefix to reuse throughout an acceptance test.
func Prefix(t *testing.T, options PrefixOptions) string {
	t.Helper()
	if options.Characters == "" {
		options.Characters = "abcdefghijklmnopqrstuvwxyz0123456789"
	}
	if options.Length == 0 {
		options.Length = 8
	}
	if options.Length < 1 {
		t.Fatal("Test prefix length must be positive")
	}
	prefix := make([]rune, options.Length)
	characters := []rune(options.Characters)
	for index := range prefix {
		choice, err := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		if err != nil {
			t.Fatalf("Generate test prefix: %v", err)
		}
		prefix[index] = characters[choice.Int64()]
	}
	value := string(prefix)
	t.Logf("Fixture prefix: %s", value)
	return value
}
