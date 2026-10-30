package services

import (
	"errors"
	"time"

	"github.com/RaucousDave/streakly/database"
	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/models"
	"gorm.io/gorm"
)

func GrantWeeklyFreeze(db *gorm.DB, userId string) error {
	var freezeLog models.FreezeLog

	now := time.Now().UTC()
	daysSinceMonday := (int(now.Weekday()+6) % 7)
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-daysSinceMonday, 0, 0, 0, 0, time.UTC)

	// err := db.Where("user_id = ? and reason = ? AND created_at >= ?", userId, models.FreezeReasonWeeklyGrant, weekStart).First(&freezeLog).Error

	database.EntDB.Query().Where(freezeLog.User)
	var newFreezeLog models.FreezeLog
	newFreezeLog.UserID = userId
	newFreezeLog.Amount = +1
	newFreezeLog.Reason = models.FreezeReasonWeeklyGrant

	error := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newFreezeLog).Error; err != nil {
			return err
		}

		res := tx.Model(&models.User{}).Where("id = ?", userId).UpdateColumn("freeze_balance", gorm.Expr("freeze_balance + ?", 1))

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})

	if error != nil {
		return error
	}

	return nil
}

func ManualFreezeHabit(db *gorm.DB, userId string, targetDate string, habitId string) (*dto.ManualFreezeResponse, error) {
	var habit models.Habit
	var freezeLog models.FreezeLog
	var freezeBalance int
	var user models.User
	err := db.Where("user_id = ? and id = ?", userId, habitId).First(&habit).Error
	if err != nil {
		return nil, err
	}

	error := db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&user).Where("id = ? and freeze_balance > 0", userId).UpdateColumn("freeze_balance", gorm.Expr("freeze_balance - ?", 1))

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.New("No freezes available")
		}

		var existing models.FreezeLog
		err := tx.Model(&existing).Where("user_id = ? and habit_id = ? and target_date = ? and reason = ?", userId, habitId, targetDate, models.FreezeReasonManualSpend).First(&existing).Error

		if err == nil {
			return errors.New("This date is already protected")
		}

		if err != gorm.ErrRecordNotFound {
			return err
		}

		freezeLog.UserID = userId
		freezeLog.Amount = -1
		freezeLog.HabitID = &habitId
		freezeLog.TargetDate = targetDate
		freezeLog.Reason = models.FreezeReasonManualSpend

		if err = tx.Create(&freezeLog).Error; err != nil {
			return err
		}

		return nil
	})

	if error != nil {
		return nil, error
	}

	err = db.Model(&user).Where("id = ?", userId).Pluck("freeze_balance", &freezeBalance).Error

	if err != nil {
		return nil, err
	}

	res := dto.ManualFreezeResponse{
		Success:       true,
		FreezeBalance: freezeBalance,
		TargetDate:    targetDate,
		Reason:        string(models.FreezeReasonManualSpend),
	}

	return &res, nil
}

func PerformanceFreezeAllocation(db *gorm.DB, userId string, weekStart time.Time, weekEnd time.Time) (*dto.PerformanceFreezeResponse, error) {
	var res dto.PerformanceFreezeResponse

	if weekStart.IsZero() || weekEnd.IsZero() {
		return nil, errors.New("Invalid date boundaries")
	}

	weekStart = weekStart.UTC()
	weekEnd = weekEnd.UTC()

	if weekEnd.Before(weekStart) || weekEnd.Equal(weekStart) {
		return nil, errors.New("invalid week range")
	}

	if weekEnd.Sub(weekStart) != 7*24*time.Hour {
		return nil, errors.New("number of days between week starting")
	}

	if weekStart.Weekday() != time.Monday || weekEnd.Weekday() != time.Monday {
		return nil, errors.New("week start and week end must be a monday")
	}

	var freezeBalance int
	if err := db.Model(&models.User{}).Where("id = ?", userId).Pluck("freeze_balance", &freezeBalance).Error; err != nil {
		return nil, err
	}
	var freezeLog models.FreezeLog
	err := db.Model(&freezeLog).Where("user_id = ? AND reason = ? AND created_at >= ? AND created_at < ?", userId, models.FreezeReasonPerformanceReward, weekStart, weekEnd).First(&freezeLog).Error

	if err == nil {
		res = dto.PerformanceFreezeResponse{
			Success:        true,
			Qualified:      true,
			FreezeGranted:  false,
			AlreadyGranted: true,
			FreezeBalance:  &freezeBalance,
		}

		return &res, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	var eligibleHabitDays int64
	var completedHabitDays int

	for current := weekStart; current.Before(weekEnd); current = current.AddDate(0, 0, 1) {
		var activeHabits int64
		var completedHabits int64
		date := current.Format("2006-01-02")
		if err := db.Where("user_id = ? AND created_at <= ? AND (is_archived = false OR archived_at >= ?)", userId, date, date).Model(&models.Habit{}).Count(&activeHabits).Error; err != nil {
			return nil, err
		}

		eligibleHabitDays += activeHabits

		if err := db.Model(&models.HabitCompletion{}).Where("user_id = ? AND is_mvt_met = true AND completed_date = ?", userId, date).Count(&completedHabits).Error; err != nil {
			return nil, err
		}

		completedHabitDays += int(completedHabits)

	}

	if eligibleHabitDays == 0 {
		return nil, errors.New("Create a habit first")
	}

	completionRate := completedHabitDays*100 >= int(eligibleHabitDays)*80

	if !completionRate {
		res = dto.PerformanceFreezeResponse{
			Success:       true,
			Qualified:     false,
			FreezeGranted: false,
			FreezeBalance: &freezeBalance,
		}

		return &res, nil
	}

	var newFreezeLog models.FreezeLog
	var updatedFreezeBalance int
	var user models.User

	newFreezeLog.UserID = userId
	newFreezeLog.Amount = +1
	newFreezeLog.Reason = models.FreezeReasonPerformanceReward

	error := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newFreezeLog).Error; err != nil {
			return err
		}

		res := tx.Model(&models.User{}).Where("id = ?", userId).UpdateColumn("freeze_balance", gorm.Expr("freeze_balance + 1"))

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})

	if error != nil {
		return nil, error
	}

	if err = db.Model(&user).Where("id = ?", userId).Pluck("freeze_balance", &updatedFreezeBalance).Error; err != nil {
		return nil, err
	}
	res = dto.PerformanceFreezeResponse{
		Success:       true,
		Qualified:     true,
		FreezeGranted: true,
		FreezeBalance: &updatedFreezeBalance,
		Reason:        string(models.FreezeReasonPerformanceReward),
	}

	return &res, nil

}
