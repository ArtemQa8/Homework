package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken создаёт подписанный JWT для указанного логина.
func GenerateToken(login string, secret []byte, ttl time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   login,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("не удалось подписать токен: %w", err)
	}
	return signed, nil
}

// ValidateToken проверяет подпись и срок действия токена.
// Возвращает true, если токен валиден.
func ValidateToken(tokenString string, secret []byte) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Проверяем, что алгоритм подписи — HMAC (HS256).
		// Без этой проверки злоумышленник мог бы подсунуть "none" или RS256.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return false, err
	}
	return token.Valid, nil
}
