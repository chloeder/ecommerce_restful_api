package utils

import (
	"math/rand"
	"time"
)

func PasscodeGenerator(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_-+="
	randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))

	random := make([]byte, length)
	for i := range random {
		random[i] = charset[randomGenerator.Intn(len(charset))]
	}

	return string(random)
}
