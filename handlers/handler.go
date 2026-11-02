package handlers

import "github.com/RaucousDave/streakly/services"

type Handler struct {
	Auth       *AuthHandler
	Habits     *HabitHandler
	Completion *CompletionHandler
	Dashboard  *DashbaordHandler
}

func NewHandlers(s *services.Services) *Handler {
	return &Handler{
		Auth:       NewAuthHandler(s.Auth),
		Habits:     NewHabitHandler(s.Habit, s.Freeze),
		Completion: NewHabitCompletionHandler(s.Habit),
		Dashboard:  NewDashboardHandler(s.Dashboard, s.User, s.Habit, s.Freeze, s.Streak),
	}
}
