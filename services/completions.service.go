package services

import (
	"context"
	"fmt"
	"time"

	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/ent"
	"github.com/RaucousDave/streakly/ent/habit"
	"github.com/RaucousDave/streakly/ent/habitcompletion"
	"github.com/google/uuid"
)

type CompletionService struct {
	db *ent.Client
}

func NewCompletionService(db *ent.Client) *CompletionService {
	return &CompletionService{db: db}
}

func (s *CompletionService) SyncCompletions(c context.Context, batch []dto.HabitCompletionDto, userId string) (int, error) {
	userID, err := uuid.Parse(userId)
	if err != nil {
		return 0, err
	}

	tx, err := s.db.Tx(c)
	if err != nil {
		return 0, nil
	}

	defer func() {
		_ = tx.Rollback()
	}()

	synced := 0

	for _, item := range batch {
		habitId, err := uuid.Parse(item.HabitId)
		if err != nil {
			return 0, fmt.Errorf("Invalid Habit ID")
		}

		completedAt, err := time.Parse(time.RFC3339, item.CompletedAtUTC)
		if err != nil {
			return 0, fmt.Errorf("Invalid date format")
		}

		existsKey, err := s.db.HabitCompletion.Query().Where(habitcompletion.IdempotencyKey(item.IdempotencyKey)).Exist(c)

		if err != nil {
			return 0, err
		}

		if existsKey {
			continue
		}

		existsDate, err := s.db.HabitCompletion.Query().Where(habitcompletion.HabitIDEQ(habitId), habitcompletion.CompletedAtUtc(completedAt)).Exist(c)

		if err != nil {
			return 0, err
		}

		if existsDate {
			continue
		}

		currentHabit, err := s.db.Habit.Query().Where(habit.UserID(userID), habit.ID(habitId)).Only(c)

		if err != nil {
			return 0, err
		}

		isMvtMet := false
		isTargetMet := false

		switch currentHabit.Type {
		case habit.TypeBOOLEAN:
			isTargetMet = true
			isMvtMet = true
		case habit.TypeQUANTITATIVE:
			if currentHabit.MvtQty != nil {
				isMvtMet = item.QtyAchieved >= *currentHabit.TargetQty
			}

			if currentHabit.TargetQty != nil {
				isTargetMet = item.QtyAchieved >= *currentHabit.TargetQty
			}
		}

		_, err = s.db.HabitCompletion.Create().SetCompletedAtUtc(completedAt).SetCompletedDate(item.CompletedDate).SetHabitID(habitId).SetUserID(userID).SetQtyAchieved(item.QtyAchieved).SetIsTargetMet(isTargetMet).SetIsMvtMet(isMvtMet).SetIdempotencyKey(item.IdempotencyKey).Save(c)

		if err != nil {
			return 0, nil
		}

		synced++
	}

	if err := tx.Commit(); err != nil {
		return 0, nil
	}

	return synced, nil
}
