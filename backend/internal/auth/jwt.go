package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTProvider interface {
	GenerateToken(userID string, ttl time.Duration) (string, error)
	ParseToken(tokenStr string) (*jwt.RegisteredClaims, error)
}

type jwtProvider struct {
	secret []byte
	issuer string
}

func NewJWTProvider(secret, issuer string) JWTProvider {
	return &jwtProvider{
		secret: []byte(secret),
		issuer: issuer,
	}
}

func (j *jwtProvider) GenerateToken(userID string, ttl time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    j.issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *jwtProvider) ParseToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
