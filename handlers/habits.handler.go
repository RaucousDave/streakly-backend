package handlers

import (
	"log"
	"strings"

	"github.com/RaucousDave/streakly/database"
	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/models"
	"github.com/RaucousDave/streakly/services"
	"github.com/RaucousDave/streakly/utils"
	"github.com/gofiber/fiber/v3"
)

func AddHabit(c fiber.Ctx) error {
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

	habit, err := services.CreateHabit(c, userId, req)
	if err != nil {
		return utils.Internal(c, err)
	}
	return utils.Created(c, "Habit created successfully", habit)
}

func GetAllHabits(c fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok || userId == "" {
		log.Print("How far relax")
		return utils.Unauthorized(c, "Unauthorized", nil)
	}

	habits, err := services.GetActiveHabits(c, userId)
	if err != nil {
		return utils.Internal(c, err)
	}

	return utils.OK(c, habits)
}

func ArchiveHabits(c fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(string)
	if !ok || userId == "" {
		return utils.Unauthorized(c, "Unauthorized", nil)
	}

	habitId := c.Params("id")

	var habit models.Habit
	message, err := services.ArchiveHabits(c, userId, habitId)
	if err != nil {
		return utils.Failure(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, nil)
	}

	return utils.Success(c, fiber.StatusOK, message, habit)
}

func FreezeHabit(c fiber.Ctx) error {
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

	res, err := services.ManualFreezeHabit(database.DB, userId, targetDate, habitId)
	if err != nil {
		return utils.Failure(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil)
	}
	return utils.Created(c, "Habit frozen successfully", res)
}
