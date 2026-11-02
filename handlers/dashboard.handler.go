package handlers

import (
	"fmt"
	"time"

	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/services"
	"github.com/RaucousDave/streakly/utils"
	"github.com/gofiber/fiber/v3"
)

type DashbaordHandler struct {
	dashboardService *services.DashboardService
	userService      *services.UserService
	habitService     *services.HabitService
	streakService    *services.StreakService
	freezeService    *services.FreezeService
}

func NewDashboardHandler(dashboardService *services.DashboardService, userService *services.UserService, habitService *services.HabitService, freezeService *services.FreezeService, streakService *services.StreakService) *DashbaordHandler {
	return &DashbaordHandler{
		dashboardService: dashboardService,
		userService:      userService,
		habitService:     habitService,
		freezeService:    freezeService,
		streakService:    streakService,
	}
}

func (h *DashbaordHandler) GetDashboard(c fiber.Ctx) error {
	const dateFormat = "2006-01-02"

	userId, ok := c.Locals("userId").(string)
	if !ok || userId == "" {
		return utils.Unauthorized(c, "Unauthorized", nil)
	}

	user, err := h.userService.GetUserStats(c, userId)
	if err != nil {
		return utils.Internal(c, err)
	}
	if user == nil {
		return utils.Internal(c, fmt.Errorf("user lookup returned nil"))
	}
	loc, err := time.LoadLocation(user.Timezone)
	if err != nil {
		return utils.Internal(c, err)
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

	heatmap, err := h.dashboardService.GetHeatmapEntry(c, from, to, userId)
	if err != nil {
		return utils.Internal(c, err)
	}

	today := time.Now().In(loc).Format("2006-01-02")

	habits, err := h.habitService.GetActiveHabits(c, userId)
	if err != nil {
		return utils.Internal(c, err)
	}

	habitItems := make([]dto.HabitDashboardItem, 0)
	if err := h.freezeService.GrantWeeklyFreeze(c, userId); err != nil {
		return utils.Internal(c, err)
	}

	for _, habit := range habits {
		todayCompletion, err := h.habitService.GetTodayCompletions(c, habit.ID.String(), today)
		if err != nil {
			return utils.Internal(c, err)
		}

		streak, err := h.streakService.CalculateStreak(c, habit.ID.String(), user.Timezone)
		if err != nil {
			return utils.Internal(c, err)
		}

		var qtyAchievedToday int
		if todayCompletion != nil {
			qtyAchievedToday = todayCompletion.QtyAchieved
		}

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
			IsMVTMetToday:    todayCompletion != nil && todayCompletion.IsMvtMet,
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
