package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"mod.go/internal/auth"
)

// AuthRequired — middleware, который проверяет JWT в заголовке Authorization.
// Если токен валиден — пропускает запрос дальше.
// Если нет — возвращает 401 и прерывает обработку.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Ошибка": "требуется заголовок Authorization"})
			return
		}

		// Ожидаем формат "Bearer <token>".
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Ошибка": "неверный формат заголовка, ожидается 'Bearer <token>'"})
			return
		}

		token := parts[1]
		valid, err := auth.ValidateToken(token, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Ошибка": "невалидный токен: " + err.Error()})
			return
		}
		if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Ошибка": "токен недействителен"})
			return
		}

		c.Next()
	}
}
