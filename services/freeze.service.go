package services

import (
	"context"
	"errors"
	"time"

	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/ent"
	"github.com/RaucousDave/streakly/ent/freezelog"
	"github.com/RaucousDave/streakly/ent/habit"
	"github.com/RaucousDave/streakly/ent/habitcompletion"
	"github.com/RaucousDave/streakly/ent/user"
	"github.com/google/uuid"
)

type FreezeService struct {
	db *ent.Client
}

func NewFreezeService(db *ent.Client) *FreezeService {
	return &FreezeService{db: db}
}

func (s *FreezeService) GrantWeeklyFreeze(c context.Context, userId string) error {

	parsedUserID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	now := time.Now()
	daysSinceMonday := (int(now.Weekday() + 6)) % 7
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-daysSinceMonday, 0, 0, 0, 0, time.UTC)

	tx, err := s.db.Tx(c)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	existing, err := tx.FreezeLog.Query().Where(freezelog.UserID(parsedUserID), freezelog.ReasonEQ(freezelog.ReasonWEEKLY_GRANT), freezelog.CreateTimeGTE(weekStart)).First(c)

	if err == nil && existing != nil {
		return tx.Commit()
	}

	if err != nil && !ent.IsNotFound(err) {
		return err
	}

	_, err = tx.FreezeLog.Create().SetUserID(parsedUserID).SetAmount(1).SetReason(freezelog.ReasonWEEKLY_GRANT).SetTargetDate(weekStart.Format("2006-01-02")).Save(c)

	if err != nil {
		return err
	}
	return nil
}

func (s *FreezeService) ManualFreezeHabit(c context.Context, userId string, targetDate string, habitId string) (*dto.ManualFreezeResponse, error) {

	parsedUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	parsedHabitUUID, err := uuid.Parse(habitId)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Tx(c)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	habitExists, err := tx.Habit.Query().Where(habit.IDEQ(parsedHabitUUID), habit.UserID(parsedUUID)).Exist(c)

	if !habitExists {
		return nil, errors.New("Habit not found")
	}

	alreadyFrozen, err := tx.FreezeLog.Query().Where(freezelog.UserID(parsedUUID), freezelog.HabitID(parsedHabitUUID), freezelog.TargetDate(targetDate), freezelog.ReasonEQ(freezelog.ReasonMANUAL_SPEND)).Exist(c)

	if err != nil {
		return nil, err
	}

	if alreadyFrozen {
		return nil, errors.New("this date is already protected")
	}

	updated, err := tx.User.Update().Where(user.ID(parsedUUID), user.FreezeBalanceGT(0)).AddFreezeBalance(-1).Save(c)

	if err != nil {
		return nil, err
	}

	if updated == 0 {
		return nil, errors.New("no freezes available")
	}

	_, err = tx.FreezeLog.Create().SetUserID(parsedUUID).SetHabitID(parsedHabitUUID).SetTargetDate(targetDate).SetAmount(-1).SetReason(freezelog.ReasonMANUAL_SPEND).Save(c)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	userRecord, err := s.db.User.Get(c, parsedUUID)
	if err != nil {
		return nil, err
	}

	return &dto.ManualFreezeResponse{
		Success:       true,
		FreezeBalance: userRecord.FreezeBalance,
		TargetDate:    targetDate,
		Reason:        freezelog.ReasonMANUAL_SPEND.String(),
	}, err
}

// func (s *FreezeService) PerformanceFreezeAllocation(db *gorm.DB, userId string, weekStart time.Time, weekEnd time.Time) (*dto.PerformanceFreezeResponse, error) {
// 	var res dto.PerformanceFreezeResponse

// 	if weekStart.IsZero() || weekEnd.IsZero() {
// 		return nil, errors.New("Invalid date boundaries")
// 	}

// 	weekStart = weekStart.UTC()
// 	weekEnd = weekEnd.UTC()

// 	if weekEnd.Before(weekStart) || weekEnd.Equal(weekStart) {
// 		return nil, errors.New("invalid week range")
// 	}

// 	if weekEnd.Sub(weekStart) != 7*24*time.Hour {
// 		return nil, errors.New("number of days between week starting")
// 	}

// 	if weekStart.Weekday() != time.Monday || weekEnd.Weekday() != time.Monday {
// 		return nil, errors.New("week start and week end must be a monday")
// 	}

// 	var freezeBalance int
// 	if err := db.Model(&models.User{}).Where("id = ?", userId).Pluck("freeze_balance", &freezeBalance).Error; err != nil {
// 		return nil, err
// 	}
// 	var freezeLog models.FreezeLog
// 	err := db.Model(&freezeLog).Where("user_id = ? AND reason = ? AND created_at >= ? AND created_at < ?", userId, models.FreezeReasonPerformanceReward, weekStart, weekEnd).First(&freezeLog).Error

// 	if err == nil {
// 		res = dto.PerformanceFreezeResponse{
// 			Success:        true,
// 			Qualified:      true,
// 			FreezeGranted:  false,
// 			AlreadyGranted: true,
// 			FreezeBalance:  &freezeBalance,
// 		}

// 		return &res, nil
// 	}
// 	if err != gorm.ErrRecordNotFound {
// 		return nil, err
// 	}

// 	var eligibleHabitDays int64
// 	var completedHabitDays int

// 	for current := weekStart; current.Before(weekEnd); current = current.AddDate(0, 0, 1) {
// 		var activeHabits int64
// 		var completedHabits int64
// 		date := current.Format("2006-01-02")
// 		if err := db.Where("user_id = ? AND created_at <= ? AND (is_archived = false OR archived_at >= ?)", userId, date, date).Model(&models.Habit{}).Count(&activeHabits).Error; err != nil {
// 			return nil, err
// 		}

// 		eligibleHabitDays += activeHabits

// 		if err := db.Model(&models.HabitCompletion{}).Where("user_id = ? AND is_mvt_met = true AND completed_date = ?", userId, date).Count(&completedHabits).Error; err != nil {
// 			return nil, err
// 		}

// 		completedHabitDays += int(completedHabits)

// 	}

// 	if eligibleHabitDays == 0 {
// 		return nil, errors.New("Create a habit first")
// 	}

// 	completionRate := completedHabitDays*100 >= int(eligibleHabitDays)*80

// 	if !completionRate {
// 		res = dto.PerformanceFreezeResponse{
// 			Success:       true,
// 			Qualified:     false,
// 			FreezeGranted: false,
// 			FreezeBalance: &freezeBalance,
// 		}

// 		return &res, nil
// 	}

// 	var newFreezeLog models.FreezeLog
// 	var updatedFreezeBalance int
// 	var user models.User

// 	newFreezeLog.UserID = userId
// 	newFreezeLog.Amount = +1
// 	newFreezeLog.Reason = models.FreezeReasonPerformanceReward

// 	error := db.Transaction(func(tx *gorm.DB) error {
// 		if err := tx.Create(&newFreezeLog).Error; err != nil {
// 			return err
// 		}

// 		res := tx.Model(&models.User{}).Where("id = ?", userId).UpdateColumn("freeze_balance", gorm.Expr("freeze_balance + 1"))

// 		if res.Error != nil {
// 			return res.Error
// 		}

// 		if res.RowsAffected == 0 {
// 			return gorm.ErrRecordNotFound
// 		}
// 		return nil
// 	})

// 	if error != nil {
// 		return nil, error
// 	}

// 	if err = db.Model(&user).Where("id = ?", userId).Pluck("freeze_balance", &updatedFreezeBalance).Error; err != nil {
// 		return nil, err
// 	}
// 	res = dto.PerformanceFreezeResponse{
// 		Success:       true,
// 		Qualified:     true,
// 		FreezeGranted: true,
// 		FreezeBalance: &updatedFreezeBalance,
// 		Reason:        string(models.FreezeReasonPerformanceReward),
// 	}

// 	return &res, nil

// }

func (s *FreezeService) PerfomanceFreezeAllocation(c context.Context, userId string, weekStart time.Time, weekEnd time.Time) (*dto.PerformanceFreezeResponse, error) {
	parsedUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	weekStart = weekStart.UTC()
	weekEnd = weekEnd.UTC()

	if weekStart.IsZero() || weekEnd.IsZero() {
		return nil, errors.New("invalid date boundaries")
	}

	if weekEnd.Sub(weekStart) != 7*24*time.Hour || weekEnd.Weekday() != time.Monday || weekStart.Weekday() != time.Monday {
		return nil, errors.New("invalid week range")
	}

	userRecord, err := s.db.User.Get(c, parsedUUID)
	if err != nil {
		return nil, err
	}

	freezeBalance := userRecord.FreezeBalance

	alreadyGranted, err := s.db.FreezeLog.Query().Where(freezelog.UserIDEQ(parsedUUID), freezelog.ReasonEQ(freezelog.ReasonPERFORMANCE_REWARD), freezelog.CreateTimeGTE(weekStart), freezelog.CreateTimeLT(weekEnd)).Exist(c)

	if err != nil {
		return nil, err
	}

	if alreadyGranted {
		return &dto.PerformanceFreezeResponse{
			Success:        true,
			Qualified:      true,
			AlreadyGranted: true,
			FreezeGranted:  false,
			FreezeBalance:  &freezeBalance,
		}, nil
	}

	var eligibleHabitDays int
	var completedHabitDays int

	for current := weekStart; current.Before(weekEnd); current = current.AddDate(0, 0, 1) {
		nextDay := current.AddDate(0, 0, 1)

		activeHabits, err := s.db.Habit.Query().Where(habit.UserIDEQ(parsedUUID), habit.CreateTimeLTE(weekStart), habit.Or(habit.ArchivedAtGTE(nextDay), habit.IsArchived(false), habit.ArchivedAtIsNil())).Count(c)

		if err != nil {
			return nil, err
		}

		completedHabits, err := s.db.HabitCompletion.Query().Where(habitcompletion.UserIDEQ(parsedUUID), habitcompletion.CompletedDateEQ(current.Format("2006-01-02")), habitcompletion.IsMvtMetEQ(true)).Count(c)

		if err != nil {
			return nil, err
		}

		eligibleHabitDays += activeHabits
		completedHabitDays += completedHabits
	}

	if eligibleHabitDays == 0 {
		return nil, errors.New("create a habit first")
	}

	completionRate := completedHabitDays * 100 / eligibleHabitDays
	qualified := completionRate >= 80

	if !qualified {
		return &dto.PerformanceFreezeResponse{
			Success:        true,
			Qualified:      false,
			CompletionRate: completionRate,
			FreezeBalance:  &freezeBalance,
		}, nil
	}

	tx, err := s.db.Tx(c)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	alreadyGranted, err = tx.FreezeLog.Query().Where(freezelog.ReasonEQ(freezelog.ReasonPERFORMANCE_REWARD), freezelog.UserID(parsedUUID), freezelog.CreateTimeGTE(weekStart), freezelog.CreateTimeLT(weekEnd)).Exist(c)

	if err != nil {
		return nil, err
	}

	if alreadyGranted {
		if err := tx.Commit(); err != nil {
			return nil, err
		}

		return &dto.PerformanceFreezeResponse{
			Success:        true,
			Qualified:      true,
			FreezeGranted:  false,
			AlreadyGranted: true,
			FreezeBalance:  &freezeBalance,
		}, nil
	}

	if _, err = tx.FreezeLog.Create().SetUserID(parsedUUID).SetAmount(1).SetReason(freezelog.ReasonPERFORMANCE_REWARD).SetTargetDate(weekStart.Format("2006-01-02")).Save(c); err != nil {
		return nil, err
	}

	updated, err := tx.User.Update().Where(user.IDEQ(parsedUUID)).AddFreezeBalance(+1).Save(c)
	if err != nil {
		return nil, err
	}

	if updated == 0 {
		return nil, errors.New("user not found")
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	freezeBalance++

	return &dto.PerformanceFreezeResponse{
		Success:        true,
		FreezeGranted:  true,
		Qualified:      true,
		CompletionRate: completionRate,
		Reason:         freezelog.ReasonPERFORMANCE_REWARD.String(),
		FreezeBalance:  &freezeBalance,
	}, nil
}
