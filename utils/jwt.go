package utils

import (
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
	"time"
)

const secret = "secret"

func Generate(userID int64, expiration time.Duration) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":     userID,
		"expiration": time.Now().Add(expiration),
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Parse(tokenString string) (*Session, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	userID := int64(claims["userid"].(float64))

	expiration, err := time.Parse(time.RFC3339, claims["expiration"].(string))
	if err != nil {
		return nil, err
	}
	res := &Session{
		UserID: userID,
		Expire: expiration,
	}
	return res, nil
}

type Session struct {
	UserID int64
	Expire time.Time
}
