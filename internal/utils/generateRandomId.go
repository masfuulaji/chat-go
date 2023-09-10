package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomID(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    seed := rand.NewSource(time.Now().UnixNano())
    random := rand.New(seed)

    id := make([]byte, length)
    for i := range id {
        id[i] = charset[random.Intn(len(charset))]
    }

    return string(id)
}
