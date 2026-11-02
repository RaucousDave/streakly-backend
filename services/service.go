package services

import "github.com/RaucousDave/streakly/ent"

type Services struct {
	Auth      *AuthService
	User      *UserService
	Habit     *HabitService
	Freeze    *FreezeService
	Dashboard *DashboardService
	Streak    *StreakService
}

func NewService(db *ent.Client) *Services {
	return &Services{
		Auth:      NewAuthService(db),
		User:      NewUserService(db),
		Habit:     NewHabitsService(db),
		Freeze:    NewFreezeService(db),
		Dashboard: NewDashboardService(db),
		Streak:    NewStreakService(db),
	}
}
