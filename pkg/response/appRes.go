package response

import "github.com/gofiber/fiber/v2"

type ResponseType struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Body    any    `json:"body"`
}

type ErrResponseType struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   any    `json:"error"`
}

func SetResponse(c *fiber.Ctx, status int, message string, body any) error {
	return c.JSON(ResponseType{
		Status:  status,
		Message: message,
		Body:    body,
	})
}

func SetErrResponse(c *fiber.Ctx, status int, message string, err any) error {
	return c.JSON(ErrResponseType{
		Status:  status,
		Message: message,
		Error:   err,
	})
}
