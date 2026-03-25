package router

import (
	"os"

	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/handler"
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/pkg/response"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func SetupRoutes(
	app 			*fiber.App,
	authHdl 		*handler.AuthHandler,
) {

	app.Get("/", func (c *fiber.Ctx) error {
		return response.SetResponse(c, fiber.StatusOK, "Server is running.", nil);
	});

	router := app.Group("/api/v1");

	authService := router.Group("/auth_service");
	authService.Post("/credential", authHdl.CredentialAuth);

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}));

	
};