package scraper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func generateHash(arr []string) string {
	// Concatenate all strings in the array
	str := ""
	for _, s := range arr {
		str += s
	}

	// Generate hash using SHA256 algorithm
	hash := sha256.Sum256([]byte(str))

	// Convert hash to string and return
	return hex.EncodeToString(hash[:])
}

func GenerateAWB(start, end int) <-chan string {
	out := make(chan string)
	go func() {
		for i := start; i <= end; i++ {
			AWBNumber := fmt.Sprintf("PRD%09d", i)
			out <- AWBNumber
		}
		close(out)
	}()
	return out
}
