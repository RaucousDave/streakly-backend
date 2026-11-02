package handlers

import (
	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/services"
	"github.com/RaucousDave/streakly/utils"
	"github.com/gofiber/fiber/v3"
)

type CompletionHandler struct {
	habitService *services.HabitService
}

func NewHabitCompletionHandler(s *services.HabitService) *CompletionHandler {
	return &CompletionHandler{habitService: s}
}
func (h *CompletionHandler) SyncHabits(c fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok || userId == "" {
		return utils.Unauthorized(c, "Unauthorized", nil)
	}

	var batch []dto.HabitCompletionDto

	if err := c.Bind().Body(&batch); err != nil {
		return utils.BadRequest(c, "Invalid request body", nil)
	}

	syncedCount, err := h.habitService.SyncCompletions(c, userId, batch)
	if err != nil {
		return utils.BadRequest(c, "Bad Request", nil)
	}

	return utils.Success(c, fiber.StatusOK, "Completions synced successfully", syncedCount)
}
