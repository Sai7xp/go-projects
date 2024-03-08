package utils

import (
	"math/rand"
	"time"
)

// generates random build Id
func GenerateBuildId() string {
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Initialize a random seed
	x1 := rand.NewSource(time.Now().UnixNano())
	y1 := rand.New(x1)

	randomString := make([]byte, 8)

	for i := range randomString {
		randomString[i] = characters[y1.Intn(len(characters))]
	}

	// Convert the byte slice to a string
	return string(randomString)
}
