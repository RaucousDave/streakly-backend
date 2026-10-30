package handlers

import (
	"errors"
	"time"

	"github.com/RaucousDave/streakly/database"
	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/models"
	"github.com/RaucousDave/streakly/utils"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SyncHabits(c fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok || userId == "" {
		return utils.Unauthorized(c, "Unauthorized", nil)
	}

	var batch []dto.HabitCompletionDto

	if err := c.Bind().Body(&batch); err != nil {
		return utils.BadRequest(c, "Invalid request body", nil)
	}

	if len(batch) == 0 {
		return utils.OK(c, nil)
	}

	syncedCount := 0

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range batch {
			var existing models.HabitCompletion

			err := tx.Where("idempotency_key = ?", item.IdempotencyKey).First(&existing).Error
			if err == nil {
				syncedCount++
				continue
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.Failure(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), nil)
			}

			var habit models.Habit

			err = tx.Where("user_id = ? and habit_id = ?", userId, item.HabitId).First(&habit).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return utils.NotFound(c, "Habit does not exist", nil)
				}
				return utils.Internal(c, err)
			}

			completed, err := time.Parse(time.RFC3339, item.CompletedAtUTC)
			if err != nil {
				completed = time.Now().UTC()
			}

			isMvtMet := false
			isTargetQtyMet := false

			switch habit.Type {
			case models.HabitTypeBoolean:
				isMvtMet = true
				isTargetQtyMet = true

			case models.HabitTypeQuantitative:
				if habit.MVTQty != nil && item.QtyAchieved >= *habit.MVTQty {
					isMvtMet = true
				}

				if habit.TargetQty != nil && item.QtyAchieved >= *habit.TargetQty {
					isTargetQtyMet = true
				}
			}

			completedHabit := models.HabitCompletion{
				HabitId:        habit.ID,
				UserID:         userId,
				CompletedDate:  item.CompletedDate,
				CompletedAtUTC: completed,
				QtyAchieved:    item.QtyAchieved,
				IsMVTMet:       isMvtMet,
				IsTargetMet:    isTargetQtyMet,
				IdempotencyKey: item.IdempotencyKey,
			}

			err = tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "habit_id"}, {Name: "completed_date"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"qty_achieved", "is_mvt_met", "is_target_met", "completed_at_utc", "idempotency_key",
				}),
			}).Create(&completedHabit).Error

			if err != nil {
				return err
			}

			syncedCount++

		}
		return nil
	})

	if err != nil {
		return utils.BadRequest(c, "Bad Request", nil)
	}

	return utils.Success(c, fiber.StatusOK, "Completions synced successfully", syncedCount)
}
