package jwt

import (
	"fmt"
	"time"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	golangjwt.RegisteredClaims
}

func GenerateToken(secret string, ttl time.Duration, userID uint, email string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: golangjwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  golangjwt.NewNumericDate(now),
			ExpiresAt: golangjwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := golangjwt.NewWithClaims(golangjwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(secret, tokenString string) (*Claims, error) {
	token, err := golangjwt.ParseWithClaims(tokenString, &Claims{}, func(token *golangjwt.Token) (any, error) {
		if _, ok := token.Method.(*golangjwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
