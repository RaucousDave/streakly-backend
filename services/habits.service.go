package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/RaucousDave/streakly/database"
	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/ent"
	"github.com/RaucousDave/streakly/ent/habit"
	"github.com/RaucousDave/streakly/ent/habitcompletion"
	"github.com/google/uuid"
)

type HabitService struct {
	db *ent.Client
}

func (s *HabitService) GetActiveHabits(c context.Context, userId string) ([]*ent.Habit, error) {
	userID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	habits, err := s.db.Habit.Query().Where(habit.UserID(userID), habit.IsArchived(false)).All(c)
	return habits, err
}

func (s *HabitService) GetTodayCompletions(c context.Context, habitId string, today string) ([]*ent.HabitCompletion, error) {

	habitID, err := uuid.Parse(habitId)
	if err != nil {
		return nil, err
	}

	completedHabits, err := database.EntDB.HabitCompletion.Query().Where(habitcompletion.HabitID(habitID), habitcompletion.CompletedDate(today)).All(c)

	return completedHabits, nil
}

// Creates a habit then based on the habit type it changes the request body coming into the service
func (s *HabitService) CreateHabit(c context.Context, userId string, input dto.HabitDto) (*ent.Habit, error) {
	id, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, errors.New("Habit title cannot be empty")
	}

	create := s.db.Habit.Create().SetUserID(id).SetTitle(title).SetIsArchived(false)

	switch input.Type {
	case dto.HabitTypeBoolean:
		create.SetType(habit.TypeBOOLEAN)
	case dto.HabitTypeQuantitative:
		if input.TargetQty == nil || input.MVTQty == nil || input.Unit == nil {
			return nil, errors.New("quantitative habits require a target quantity mvt quantity and a unit")
		}
		create.SetType(habit.TypeQUANTITATIVE).SetMvtQty(*input.MVTQty).SetTargetQty(*input.TargetQty).SetUnit(*input.Unit)

	default:
		return nil, errors.New("Invalid habit type")
	}

	return create.Save(c)
}

func (s *HabitService) ArchiveHabits(c context.Context, userId string, habitId string) (string, error) {

	userID, err := uuid.Parse(userId)
	if err != nil {
		return "", err
	}

	habitID, err := uuid.Parse(habitId)
	if err != nil {
		return "", err
	}
	currentHabit, err := s.db.Habit.Query().Where(habit.UserID(userID), habit.ID(habitID)).Only(c)
	if err != nil {
		return "", err
	}

	updatedHabit := s.db.Habit.UpdateOneID(currentHabit.ID).SetIsArchived(!currentHabit.IsArchived)

	if currentHabit.IsArchived {
		updatedHabit.ClearArchivedAt()
	} else {
		updatedHabit.SetArchivedAt(time.Now())
	}

	if _, err := updatedHabit.Save(c); err != nil {
		return "", err
	}

	if currentHabit.IsArchived {
		return "Habit unarchived successfully", nil
	}

	message := "Habit archived successfully"

	return message, nil
}
