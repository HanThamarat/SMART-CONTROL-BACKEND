package encrypt

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	UserId uint `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status bool `json:"status"`
	jwt.StandardClaims
}

func JWTDecryption(c *fiber.Ctx) (*UserClaims, error) {
	// 1. Get the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header missing")
	}

	// 2. Remove "Bearer " prefix
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return nil, errors.New("authorization header format must be Bearer {token}")
	}

	// 3. Parse and Validate
	// We pass a pointer to UserClaims here
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// Handle parsing errors (including expired tokens)
	if err != nil {
		return nil, err
	}

	// 4. Extract Claims
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}