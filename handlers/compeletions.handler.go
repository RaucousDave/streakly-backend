package handlers

import (
	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/services"
	"github.com/RaucousDave/streakly/utils"
	"github.com/gofiber/fiber/v3"
)

type CompletionHandler struct {
	completionService *services.CompletionService
}

func NewCompletionHandler(completionService *services.CompletionService) *CompletionHandler {
	return &CompletionHandler{
		completionService: completionService,
	}
}
func (h *CompletionHandler) SyncCompletions(c fiber.Ctx) error {
	var req []dto.HabitCompletionDto
	userId, ok := c.Locals("userId").(string)

	if !ok || userId == "" {
		return utils.Unauthorized(c, "Unauthorized", nil)
	}

	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body", nil)
	}

	synced, err := h.completionService.SyncCompletions(c, req, userId)
	if err != nil {
		return err
	}

	return utils.Success(c, fiber.StatusOK, "Habits synced successfully", synced)
}
