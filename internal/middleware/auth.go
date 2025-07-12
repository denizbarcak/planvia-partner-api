package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// AuthMiddleware JWT token'ı doğrular ve partner ID'yi context'e ekler
func AuthMiddleware(c *fiber.Ctx) error {
	// Authorization header'ı kontrol et
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header eksik",
		})
	}

	// Bearer token'ı ayır
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Geçersiz token formatı",
		})
	}

	// Token'ı doğrula
	token, err := jwt.Parse(tokenParts[1], func(token *jwt.Token) (interface{}, error) {
		// Token'ın signing method'unu kontrol et
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Geçersiz token")
		}
		// TODO: JWT secret key'i config'den al
		return []byte("your-secret-key"), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token doğrulanamadı",
		})
	}

	// Token claims'i kontrol et
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Partner ID'yi context'e ekle
		c.Locals("partnerId", claims["partnerId"])
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Geçersiz token",
	})
} 