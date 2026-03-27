package handler

import (
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/domain"
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	usecase domain.AuthUsecase
}

func NewAuthHandler(uc domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		usecase: uc,
	}
}

func (h *AuthHandler) CredentialAuth(c *fiber.Ctx) error {
	var req domain.AuthDTO
	if err := c.BodyParser(&req); err != nil {
		return response.SetErrResponse(
			c,
			fiber.StatusBadRequest,
			"Authentication failed.",
			err.Error(),
		)
	}

	result, err := h.usecase.CredentialAuth(req)

	if err != nil {
		return response.SetResponse(
			c,
			fiber.StatusBadRequest,
			"Authentication failed.",
			err.Error(),
		)
	}

	return response.SetResponse(
		c,
		fiber.StatusOK,
		"Authentication successfully.",
		result,
	)
}
