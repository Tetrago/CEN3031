package util

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenContents struct {
	Ident string `json:"ident"`
}

var Secret []byte

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

func Hash(value string) string {
	sh := sha256.New()
	sh.Write([]byte(value))
	return fmt.Sprintf("%x", sh.Sum(nil))
}

func MakeToken(contents TokenContents) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"contents": contents,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	})
	return token.SignedString(Secret)
}

func ProcessToken(str string) (*TokenContents, error) {
	token, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return Secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if claims["exp"].(int64) < time.Now().Unix() {
			return nil, fmt.Errorf("token expired")
		}

		return &TokenContents{
			Ident: claims["contents"].(map[string]interface{})["ident"].(string),
		}, nil
	} else {
		return nil, err
	}
}
