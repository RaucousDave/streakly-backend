package models

import (
	"time"

	"gorm.io/gorm"
)

type HabitType string
type FreezeReason string

const (
	HabitTypeBoolean      HabitType = "BOOLEAN"
	HabitTypeQuantitative HabitType = "QUANTITATIVE"
)

const (
	FreezeReasonWeeklyGrant       FreezeReason = "WEEKLY_GRANT"
	FreezeReasonPerformanceReward FreezeReason = "PERFORMANCE_REWARD"
	FreezeReasonManualSpend       FreezeReason = "MANUAL_SPEND"
)

type AuthProvider string

const (
	AuthProviderGoogle AuthProvider = "GOOGLE"
	AuthProviderEmail  AuthProvider = "EMAIL"
)

type User struct {
	ID            string       `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email         string       `gorm:"column:email;uniqueIndex;not null" json:"email"`
	GoogleID      *string      `gorm:"column:google_id;uniqueIndex;not null" json:"-"`
	PasswordHash  *string      `gorm:"column:password_hash" json:"-"`
	Provider      AuthProvider `gorm:"column:auth_provider;not null" json:"auth_provider"`
	DisplayName   string       `gorm:"column:display_name;not null" json:"display_name"`
	TimeZone      string       `gorm:"column:timezone;not null" json:"time_zone"`
	AvatarUrl     string       `gorm:"column:avatar_url;not null" json:"avatar_url"`
	FreezeBalance int          `gorm:"column:freeze_balance;not null;default:1" json:"freeze_balance"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Habit struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID     string    `gorm:"index;not null" json:"user_id"`
	User       User      `gorm:"constraint:OnDelete:CASCADE;"`
	Title      string    `gorm:"not null" json:"title"`
	Type       HabitType `gorm:"type:varchar(20);not null; default:'BOOLEAN'" json:"type"`
	TargetQty  *int      `gorm:"default:null" json:"target_qty"`
	MVTQty     *int      `gorm:"default:null" json:"mvtqty"`
	Unit       *string   `gorm:"default:null" json:"unit"`
	IsArchived bool      `gorm:"default:false" json:"is_archived"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ArchivedAt time.Time      `gorm:"index" json:"archived_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type HabitCompletion struct {
	ID             string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	HabitId        string    `gorm:"index;not null" json:"habit_id"`
	Habit          Habit     `gorm:"constraint:OnDelete:CASCADE"`
	UserID         string    `gorm:"index;not null" json:"user_id"`
	CompletedDate  string    `gorm:"index:idx_habit_date,unique;type:varchar(10);not null" json:"completed_date"`
	CompletedAtUTC time.Time `gorm:"not null" json:"completed_at_utc"`
	QtyAchieved    int       `gorm:"not null" json:"qty_achieved"`
	IsMVTMet       bool      `gorm:"not null; default:false" json:"is_mvt_met"`
	IsTargetMet    bool      `gorm:"not null; default:false" json:"is_target_met"`
	IdempotencyKey string    `gorm:"uniqueIndex;not null" json:"idempotency_key"`
	CreatedAt      time.Time
}

type FreezeLog struct {
	ID         string       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Reason     FreezeReason `gorm:"not null" json:"reason"`
	HabitID    *string      `gorm:"default:null" json:"habit_id"`
	Habit      Habit        `gorm:"constraint:OnDelete:CASCADE"`
	Amount     int          `gorm:"not null" json:"amount"`
	UserID     string       `gorm:"index; not null" json:"user_id"`
	User       User         `gorm:"constraint:OnDelete:CASCADE"`
	TargetDate string       `gorm:"type:varchar(10)" json:"target_date"`
	CreatedAt  time.Time
}
