package handler

import (
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/domain"
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type WidgetHandler struct {
	usecase		domain.WidgetUsecase
}

func NewWidgetHandler(u domain.WidgetUsecase) *WidgetHandler {
	return &WidgetHandler{
		usecase: u,
	}
}

func (h *WidgetHandler) CreateNewWidget(c *fiber.Ctx) error {
	var req domain.WidgetDTO;
	if err := c.BodyParser(&req); err != nil {
		return response.SetErrResponse(
			c,
			fiber.StatusBadRequest,
			"Create a new widget failed.",
			err,
		);
	}

	result, err := h.usecase.CreateNewWidget(req);

	if err != nil {
		return response.SetErrResponse(
			c,
			fiber.StatusBadRequest,
			"Create a new widget failed.",
			err,
		);
	}

	return response.SetResponse(
		c,
		fiber.StatusCreated,
		"Create a new widget successfully.",
		result,
	);
}

func (h *WidgetHandler) FindAllWidget(c *fiber.Ctx) error {
	result, err := h.usecase.FindAllWidget();

	if err != nil {
		return response.SetErrResponse(
			c,
			fiber.StatusBadRequest,
			"Find all widgets failed.",
			err,
		);
	}

	return response.SetResponse(
		c,
		fiber.StatusOK,
		"Find all widgets successfully.",
		result,
	); 
}

func (h *WidgetHandler) FindWidget(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id");
	if err != nil {
		return response.SetErrResponse(
			c,
			fiber.StatusBadRequest,
			"Find widget by id failed.",
			err,
		);
	}

	result, err := h.usecase.FindWidget(id);

	if err != nil {
		return response.SetErrResponse(
			c,
			fiber.StatusBadRequest,
			"Find widget by id failed.",
			err,
		);
	}

	return response.SetResponse(
		c,
		fiber.StatusOK,
		"Find widget by id successfully.",
		result,
	); 
}

func (h *WidgetHandler) UpdateWidget(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id");
	if err != nil {
		return response.SetErrResponse(
			c,
			fiber.StatusBadRequest,
			"update widget failed.",
			err,
		);
	}

	var req domain.WidgetUpdateDTO;
	if err := c.BodyParser(&req); err != nil {
		return response.SetErrResponse(
			c,
			fiber.StatusBadRequest,
			"update widget failed.",
			err,
		);
	}

	result, err := h.usecase.UpdateWidget(id, req);

	if err != nil {
		return response.SetErrResponse(
			c,
			fiber.StatusBadRequest,
			"update widget failed.",
			err,
		);
	}

	return response.SetResponse(
		c,
		fiber.StatusOK,
		"Widget updated.",
		result,
	); 
}

func (h *WidgetHandler) DeleteWidget(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id");
	if err != nil {
		return response.SetErrResponse(
			c,
			fiber.StatusBadRequest,
			"Delete failed.",
			err,
		);
	}

	result, err := h.usecase.DeleteWidget(id);

	if err != nil {
		return response.SetErrResponse(
			c,
			fiber.StatusBadRequest,
			"Delete failed.",
			err,
		);
	}

	return response.SetResponse(
		c,
		fiber.StatusOK,
		"Widget deleted.",
		result,
	); 
}

