package services

import (
	"time"

	"github.com/RaucousDave/streakly/models"
	"gorm.io/gorm"
)

func CalculateStreak(db *gorm.DB, habitId string, userTimezone string, completedHabit models.HabitCompletion) (int, error) {
	var dates []string
	if err := db.Model(&models.HabitCompletion{}).Where("habit_id = ? and is_mvt_met = ?", habitId, true).Pluck("completed_date", &dates).Error; err != nil {
		return 0, err
	}

	completedDates := make(map[string]bool)

	for _, d := range dates {
		completedDates[d] = true
	}

	loc, err := time.LoadLocation(userTimezone)
	if err != nil {
		loc = time.UTC
	}

	now := time.Now().In(loc)
	cursor := now

	streak := 0
	todayStr := cursor.Format("2006-01-02")

	if completedDates[todayStr] {
		streak = 1
	}

	cursor = cursor.AddDate(0, 0, -1)

	for {
		dateStr := cursor.Format("2006-01-02")

		if completedDates[dateStr] {
			streak++
			cursor = cursor.AddDate(0, 0, -1)
		} else {
			break
		}
	}
	return streak, nil
}
