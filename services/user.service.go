package services

import (
	"context"

	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/ent"
	"github.com/google/uuid"
)

type UserService struct {
	db *ent.Client
}

func NewUserService(db *ent.Client) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUserStats(ctx context.Context, userId string) (*dto.DashboardUser, error) {
	var userInfo dto.DashboardUser

	userID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	userRecord, err := s.db.User.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &dto.DashboardUser{
          DisplayName:   userRecord.DisplayName,
          Timezone:      userRecord.Timezone,
          FreezeBalance: userRecord.FreezeBalance,
      }, nil
}
