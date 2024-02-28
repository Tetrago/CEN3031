package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tetrago/motmot/api/internal/globals"
)

type Token jwt.MapClaims

func NewToken() *Token {
	t := Token(jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})

	return &t
}

func (t Token) UserIdentifier() string {
	return t["ident"].(string)
}

func (t *Token) SetUserIdentifier(ident string) *Token {
	(*t)["ident"] = ident
	return t
}

func (t Token) Serialize() (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(t)).SignedString([]byte(globals.Opts.TokenSecret))
}

func ParseToken(raw string) (*Token, error) {
	token, err := jwt.Parse(raw, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(globals.Opts.TokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if int64(claims["exp"].(float64)) < time.Now().Unix() {
			return nil, fmt.Errorf("token expired")
		}

		t := Token(claims)
		return &t, nil
	} else {
		return nil, err
	}
}
