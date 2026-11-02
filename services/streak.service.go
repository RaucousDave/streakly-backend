package services

import (
	"context"
	"time"

	"github.com/RaucousDave/streakly/ent"
	"github.com/RaucousDave/streakly/ent/habitcompletion"
	"github.com/google/uuid"
)

type StreakService struct {
	db *ent.Client
}

func NewStreakService(db *ent.Client) *StreakService {
	return &StreakService{db: db}
}

func (s *StreakService) CalculateStreak(c context.Context, habitId string, userTimezone string) (int, error) {
	parsedHabitID, err := uuid.Parse(habitId)
	if err != nil {
		return 0, err
	}

	dates, err := s.db.HabitCompletion.Query().Where(habitcompletion.HabitIDEQ(parsedHabitID), habitcompletion.IsMvtMet(true)).Select(habitcompletion.FieldCompletedDate).Strings(c)

	if err != nil {
		return 0, err
	}

	completedDates := make(map[string]bool)

	for _, d := range dates {
		completedDates[d] = true
	}

	loc, err := time.LoadLocation(userTimezone)

	if err != nil {
		return 0, err
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
		dateStr := todayStr

		if completedDates[dateStr] {
			streak++
			cursor = cursor.AddDate(0, 0, -1)
		} else {
			break
		}
	}
	return streak, nil
}
