package services

import (
	"time"

	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/models"
	"gorm.io/gorm"
)

func GetHeatmapEntry(db *gorm.DB, from time.Time, to time.Time, userId string) (*[]dto.HeatmapEntry, error) {
	const dateFormat = "2006-01-02"
	var habits []models.Habit
	var heatmapEntries []dto.HeatmapEntry
	var completedHabit []models.HabitCompletion
	for current := from; !current.After(to); current.AddDate(0, 0, 1) {
		date := current.Format(dateFormat)
		var totalHabits int
		err := db.Where("user_id = ? and created_at <= ?  and archived_at = ?  or archived_at > ?", userId, to, "null", from).Find(&habits)

		if err != nil {
			return nil, err.Error
		}

		totalHabits = len(habits)

		var completedHabits int
		err = db.Where("user_id ? and completed_at = ?", userId, date).Find(&completedHabit)
		completedHabits = len(completedHabit)

		completed := totalHabits == completedHabits

		heatmapEntry := dto.HeatmapEntry{
			Date:           date,
			Completed:      completed,
			CompletedCount: completedHabits,
			Total:          totalHabits,
		}

		heatmapEntries = append(heatmapEntries, heatmapEntry)

	}
	return &heatmapEntries, nil
}
