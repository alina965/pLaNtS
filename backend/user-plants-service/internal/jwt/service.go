package jwt

import (
	"fmt"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret []byte
}

type AccessClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwtlib.RegisteredClaims
}

func NewService(secret []byte) *Service {
	return &Service{secret: secret}
}

func (s *Service) Parse(accessToken string) (userID string, email string, err error) {
	token, err := jwtlib.ParseWithClaims(accessToken, &AccessClaims{}, func(token *jwtlib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return "", "", fmt.Errorf("invalid token claims")
	}

	return claims.UserID, claims.Email, nil
}
