package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenContents struct {
	Ident string `json:"ident"`
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
