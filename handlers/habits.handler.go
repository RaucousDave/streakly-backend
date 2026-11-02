package handlers

import (
	"log"
	"strings"

	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/models"
	"github.com/RaucousDave/streakly/services"
	"github.com/RaucousDave/streakly/utils"
	"github.com/gofiber/fiber/v3"
)

type HabitHandler struct {
	habitService  *services.HabitService
	freezeService *services.FreezeService
}

func NewHabitHandler(habitService *services.HabitService, freezeSevice *services.FreezeService) *HabitHandler {
	return &HabitHandler{
		habitService:  habitService,
		freezeService: freezeSevice,
	}
}

func (h *HabitHandler) AddHabit(c fiber.Ctx) error {
	var req dto.HabitDto

	userId, ok := c.Locals("userId").(string)
	if !ok || userId == "" {
		return utils.Unauthorized(c, "Unauthorized", nil)
	}

	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body", nil)
	}

	if strings.TrimSpace(req.Title) == "" {
		return utils.BadRequest(c, "Habit title cannot be empty", nil)
	}

	habit, err := h.habitService.CreateHabit(c, userId, req)
	if err != nil {
		return utils.Internal(c, err)
	}
	return utils.Created(c, "Habit created successfully", habit)
}

func (h *HabitHandler) GetAllHabits(c fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok || userId == "" {
		log.Print("How far relax")
		return utils.Unauthorized(c, "Unauthorized", nil)
	}

	habits, err := h.habitService.GetActiveHabits(c, userId)
	if err != nil {
		return utils.Internal(c, err)
	}

	return utils.OK(c, habits)
}

func (h *HabitHandler) ArchiveHabits(c fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok || userId == "" {
		return utils.Unauthorized(c, "Unauthorized", nil)
	}

	habitId := c.Params("id")

	var habit models.Habit
	message, err := h.habitService.ArchiveHabits(c, userId, habitId)
	if err != nil {
		return utils.Failure(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, nil)
	}

	return utils.Success(c, fiber.StatusOK, message, habit)
}

func (h *HabitHandler) FreezeHabit(c fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(string)

	if !ok || userId == "" {
		return utils.Unauthorized(c, "Unauthorized", nil)
	}
	habitId := c.Params("id")
	var targetDate string
	err := c.Bind().Body(&targetDate)
	if err != nil {
		return utils.BadRequest(c, "Bad request", err)
	}

	res, err := h.freezeService.ManualFreezeHabit(c, userId, targetDate, habitId)
	if err != nil {
		return utils.Failure(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil)
	}
	return utils.Created(c, "Habit frozen successfully", res)
}
