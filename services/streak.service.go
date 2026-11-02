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

	cursor := time.Now().In(loc)
	streak := 0

	for {
		dateStr := cursor.Format("2006-01-02")
		if !completedDates[dateStr] {
			break
		}
		streak++
		cursor = cursor.AddDate(0, 0, -1)
	}
	return streak, nil
}
