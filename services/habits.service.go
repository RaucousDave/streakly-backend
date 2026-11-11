package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/ent"
	"github.com/RaucousDave/streakly/ent/habit"
	"github.com/RaucousDave/streakly/ent/habitcompletion"
	"github.com/google/uuid"
)

type HabitService struct {
	db *ent.Client
}

func NewHabitsService(db *ent.Client) *HabitService {
	return &HabitService{db: db}
}

func (s *HabitService) GetActiveHabits(c context.Context, userId string) ([]*ent.Habit, error) {
	userID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	habits, err := s.db.Habit.Query().Where(habit.UserID(userID), habit.IsArchived(false)).All(c)
	return habits, err
}

func (s *HabitService) GetArchivedHabits(c context.Context, userId string) ([]*ent.Habit, error) {
	userID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	return s.db.Habit.Query().Where(habit.UserID(userID), habit.IsArchived(true)).Order(ent.Desc(habit.FieldArchivedAt)).All(c)
}

func (s *HabitService) GetTodayCompletions(c context.Context, habitId string, today string) (*ent.HabitCompletion, error) {

	habitID, err := uuid.Parse(habitId)
	if err != nil {
		return nil, err
	}

	completedHabits, err := s.db.HabitCompletion.Query().Where(habitcompletion.HabitID(habitID), habitcompletion.CompletedDate(today)).Only(c)

	if ent.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

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

var ErrInvalidHabitEdit = errors.New("invalid habit edit")

func (s *HabitService) EditHabit(c context.Context, userId, habitId string, input dto.EditHabitDto) (*ent.Habit, error) {
	userID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid user ID", ErrInvalidHabitEdit)
	}
	habitID, err := uuid.Parse(habitId)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid habit ID", ErrInvalidHabitEdit)
	}
	current, err := s.db.Habit.Query().Where(habit.ID(habitID), habit.UserID(userID)).Only(c)
	if err != nil {
		return nil, err
	}
	next, err := prepareHabitEdit(*current, input)
	if err != nil {
		return nil, err
	}
	update := s.db.Habit.UpdateOneID(habitID).Where(habit.UserID(userID)).
		SetTitle(next.Title).SetType(next.Type)
	if next.Type == habit.TypeBOOLEAN {
		update.ClearUnit().ClearTargetQty().ClearMvtQty()
	} else {
		update.SetUnit(*next.Unit).SetTargetQty(*next.TargetQty).SetMvtQty(*next.MvtQty)
	}
	return update.Save(c)
}

// Validate the merged state so partial edits cannot leave inconsistent quantities.
func prepareHabitEdit(current ent.Habit, input dto.EditHabitDto) (ent.Habit, error) {
	if input.Title == nil && input.Type == nil && input.Unit == nil && input.TargetQty == nil && input.MVTQty == nil {
		return current, fmt.Errorf("%w: provide at least one field", ErrInvalidHabitEdit)
	}
	if input.Title != nil {
		current.Title = strings.TrimSpace(*input.Title)
	}
	if strings.TrimSpace(current.Title) == "" {
		return current, fmt.Errorf("%w: title cannot be empty", ErrInvalidHabitEdit)
	}
	if input.Type != nil {
		current.Type = habit.Type(*input.Type)
	}
	if input.Unit != nil {
		unit := strings.TrimSpace(*input.Unit)
		current.Unit = &unit
	}
	if input.TargetQty != nil {
		current.TargetQty = input.TargetQty
	}
	if input.MVTQty != nil {
		current.MvtQty = input.MVTQty
	}
	switch current.Type {
	case habit.TypeBOOLEAN:
		current.Unit, current.TargetQty, current.MvtQty = nil, nil, nil
	case habit.TypeQUANTITATIVE:
		if current.Unit == nil || strings.TrimSpace(*current.Unit) == "" || current.TargetQty == nil || current.MvtQty == nil {
			return current, fmt.Errorf("%w: quantitative habits require unit, target_qty and mvt_qty", ErrInvalidHabitEdit)
		}
		if *current.TargetQty <= 0 || *current.MvtQty <= 0 || *current.MvtQty > *current.TargetQty {
			return current, fmt.Errorf("%w: quantities must satisfy 0 < mvt_qty <= target_qty", ErrInvalidHabitEdit)
		}
	default:
		return current, fmt.Errorf("%w: type must be BOOLEAN or QUANTITATIVE", ErrInvalidHabitEdit)
	}
	return current, nil
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

func (s *HabitService) SyncCompletions(c context.Context, userId string, batch []dto.HabitCompletionDto) (int, error) {
	userID, err := uuid.Parse(userId)
	if err != nil {
		return 0, err
	}

	tx, err := s.db.Tx(c)

	if err != nil {
		return 0, err
	}

	defer func() {
		tx.Rollback()
	}()

	synced := 0

	for _, item := range batch {
		parsedHabitUUID, err := uuid.Parse(item.HabitId)
		if err != nil {
			return 0, err
		}

		completedAt, err := time.Parse(time.RFC3339, item.CompletedAtUTC)
		if err != nil {
			return 0, fmt.Errorf("invalid completed_at_utc for habit %s: %w", item.HabitId, err)
		}

		exists, err := tx.HabitCompletion.Query().Where(habitcompletion.IdempotencyKeyEQ(item.IdempotencyKey)).Exist(c)
		if err != nil {
			return 0, err
		}

		if exists {
			continue
		}

		currentHabit, err := tx.Habit.Query().Where(habit.UserIDEQ(userID), habit.IDEQ(parsedHabitUUID)).Only(c)

		if err != nil {
			if ent.IsNotFound(err) {
				return 0, fmt.Errorf("habit %s does not belong to the user", item.HabitId)
			}
			return 0, err
		}

		isMVTMet := false
		isTargetMet := false

		switch currentHabit.Type {
		case habit.TypeBOOLEAN:
			isMVTMet = true
			isTargetMet = true

		case habit.TypeQUANTITATIVE:
			if currentHabit.MvtQty != nil {
				isMVTMet = item.QtyAchieved >= *currentHabit.MvtQty
			}
			if currentHabit.TargetQty != nil {
				isTargetMet = item.QtyAchieved >= *currentHabit.TargetQty
			}
		}

		if _, err = tx.HabitCompletion.Create().SetHabitID(parsedHabitUUID).SetUserID(userID).SetCompletedDate(item.CompletedDate).SetCompletedAtUtc(completedAt).SetIsMvtMet(isMVTMet).SetIsTargetMet(isTargetMet).SetIdempotencyKey(item.IdempotencyKey).SetQtyAchieved(item.QtyAchieved).Save(c); err != nil {
			return 0, err
		}

		synced++
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return synced, nil
}
