package handlers

import (
	"time"

	"github.com/RaucousDave/streakly/database"
	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/services"
	"github.com/RaucousDave/streakly/utils"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func GetDashboard(c fiber.Ctx) error {
	var DB *gorm.DB
	const dateFormat = "2006-01-02"

	DB = database.DB
	userId, ok := c.Locals("userId").(string)
	if !ok || userId == "" {
		return utils.Unauthorized(c, "Unauthorized", nil)
	}

	user, err := services.GetUserStats(c, userId)
	loc, err := time.LoadLocation(user.Timezone)
	if err != nil {
		return utils.Internal(c)
	}

	todayDate := time.Now().In(loc)
	defaultTo := todayDate
	defaultFrom := defaultTo.AddDate(0, 0, -364)

	fromParam := c.Query("from")
	toParam := c.Query("to")

	from := defaultFrom
	to := defaultTo

	if fromParam != "" {

		from, err = time.ParseInLocation("2006-01-02", fromParam, loc)
		if err != nil {
			return utils.BadRequest(c, "Bad request", nil)
		}
	}

	if toParam != "" {
		to, err = time.ParseInLocation("2006-01-02", toParam, loc)
		if err != nil {
			return utils.BadRequest(c, "Bad request", nil)

		}
	}

	heatmap, err := services.GetHeatmapEntry(DB, from, to, userId)
	if err != nil {
		return utils.Internal(c, err)
	}

	today := time.Now().In(loc).Format("2006-01-02")

	habits, err := services.GetActiveHabits(c, userId)
	if err != nil {
		return utils.Internal(c, err)
	}

	habitItems := make([]dto.HabitDashboardItem, 0)

	for _, habit := range habits {
		todayCompletion, err := services.GetTodayCompletions(DB, habit.ID.String(), today)
		if err != nil {
			return utils.Internal(c, err)
		}

		streak, err := services.CalculateStreak(DB, habit.ID.String(), user.Timezone, *todayCompletion)
		if err != nil {
			return utils.Internal(c, err)
		}

		if err := services.GrantWeeklyFreeze(DB, userId); err != nil {
			return utils.Internal(c, err)
		}

		var qtyAchievedToday int
		if todayCompletion != nil {
			qtyAchievedToday = todayCompletion.QtyAchieved
		}
		qtyAchievedToday = 0

		item := dto.HabitDashboardItem{
			ID:               habit.ID.String(),
			Title:            habit.Title,
			Type:             string(habit.Type),
			TargetQty:        habit.TargetQty,
			MVTQty:           habit.MvtQty,
			Unit:             habit.Unit,
			CurrentStreak:    streak,
			IsCompletedToday: todayCompletion != nil,
			QtyAchievedToday: qtyAchievedToday,
			IsMVTMetToday:    todayCompletion != nil && todayCompletion.IsMVTMet,
			IsTargetMetToday: todayCompletion != nil && todayCompletion.IsTargetMet,
		}

		habitItems = append(habitItems, item)
	}

	response := dto.DashboardResponse{
		User: dto.DashboardUser{
			DisplayName:   user.DisplayName,
			Timezone:      user.Timezone,
			FreezeBalance: user.FreezeBalance,
		},
		Today: dto.TodaySection{
			Date:   today,
			Habits: habitItems,
		},
		Heatmap: *heatmap,
	}

	return utils.OK(c, response)
}
