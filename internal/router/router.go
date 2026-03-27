package router

import (
	"os"

	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/handler"
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/socket"
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/pkg/response"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func SetupRoutes(
	app *fiber.App,
	authHdl *handler.AuthHandler,
	widgetHdl *handler.WidgetHandler,
	socketServer *socket.Server,
) {
	app.Get("/", func(c *fiber.Ctx) error {
		return response.SetResponse(c, fiber.StatusOK, "Server is running.", nil)
	})

	app.Use("/socket.io", socketServer.Upgrade)
	app.Get("/socket.io", socketServer.Handler())

	router := app.Group("/api/v1")

	authService := router.Group("/auth_service")
	authService.Post("/credential", authHdl.CredentialAuth)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	widgetService := router.Group("/widget_service")
	widgetService.Post("/create", widgetHdl.CreateNewWidget)
	widgetService.Get("/finds", widgetHdl.FindAllWidget)
	widgetService.Get("/find/:id", widgetHdl.FindWidget)
	widgetService.Put("/update/:id", widgetHdl.UpdateWidget)
	widgetService.Delete("/delete/:id", widgetHdl.DeleteWidget)
}
