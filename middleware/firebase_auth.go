package middleware

import (
	"context"
	"github.com/gofiber/fiber/v2"

	"ecoApi/firebase"
)

func FirebaseAuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header missing",
		})
	}

	idToken := authHeader[len("Bearer "):]

	token, err := firebase.FirebaseAuth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	} 	
	c.Locals("uid", token.UID)
	return c.Next()
}