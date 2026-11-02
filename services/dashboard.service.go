package services

import (
	"context"
	"errors"
	"time"

	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/ent"
	"github.com/RaucousDave/streakly/ent/habit"
	"github.com/RaucousDave/streakly/ent/habitcompletion"
	"github.com/google/uuid"
)

type DashboardService struct {
	db *ent.Client
}

func NewDashboardService(db *ent.Client) *DashboardService {
	return &DashboardService{db: db}
}

func (s *DashboardService) GetHeatmapEntry(c context.Context, from time.Time, to time.Time, userId string) (*[]dto.HeatmapEntry, error) {

	from = from.UTC()
	to = to.UTC()

	if to.Before(from) {
		return nil, errors.New("invalid date range")
	}
	parsedUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	entries := make([]dto.HeatmapEntry, 0)

	for current := from; !current.After(to); current = current.AddDate(0, 0, 1) {
		date := current.Format("2006-01-02")
		nextDay := current.AddDate(0, 0, 1)

		totalHabits, err := s.db.Habit.Query().Where(habit.UserIDEQ(parsedUUID), habit.CreateTimeLT(nextDay), habit.Or(habit.IsArchivedEQ(false), habit.ArchivedAtIsNil(), habit.ArchivedAtGT(current))).Count(c)

		if err != nil {
			return nil, err
		}

		completedHabits, err := s.db.HabitCompletion.Query().Where(habitcompletion.UserIDEQ(parsedUUID), habitcompletion.CompletedDateEQ(date), habitcompletion.IsMvtMet(true)).Count(c)

		if err != nil {
			return nil, err
		}

		entries = append(entries, dto.HeatmapEntry{
			Date:           date,
			Completed:      totalHabits > 0 && completedHabits == totalHabits,
			Total:          totalHabits,
			CompletedCount: completedHabits,
		})

	}

	return &entries, nil
}
