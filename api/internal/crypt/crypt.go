package crypt

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

func GenerateBase64(length int) (string, error) {
	const pattern = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	val := make([]byte, length)

	for i := 0; i < length; i++ {
		r, err := rand.Int(rand.Reader, big.NewInt(int64(len(pattern))))

		if err != nil {
			return "", err
		}

		val[i] = pattern[r.Int64()]
	}

	return string(val), nil
}

func HashToInt(value string) int {
	sum := 0

	for i := 0; i < len(value); i++ {
		sum += int(value[i])
	}

	return sum
}

func Hash(value string) string {
	sh := sha256.New()
	sh.Write([]byte(value))
	return fmt.Sprintf("%x", sh.Sum(nil))
}
