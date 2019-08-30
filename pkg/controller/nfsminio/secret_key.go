package nfsminio

import (
	"math/rand"
	"strings"
	"time"
)

// Generate a random secret key from a-zA-Z0-9. The random seed is NOT cryptographically secure.
func randomSecretKey(length int) string {
	// Adapted from: https://yourbasic.org/golang/generate-random-string/
	rand.Seed(time.Now().UnixNano())  // TODO - replace with cryptographic secure random seed
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
