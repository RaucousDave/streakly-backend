package services

import (
	"context"

	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/ent"
	"github.com/RaucousDave/streakly/ent/user"
	"github.com/google/uuid"
)

type UserService struct{
	db *ent.Client
}
func (s *UserService) GetUserStats(ctx context.Context, userId string) (*dto.DashboardUser, error) {
	var userInfo dto.DashboardUser

	userID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	err = s.db.User.Query().Where(user.ID(userID)).Select(user.FieldDisplayName, user.FieldTimezone, user.FieldFreezeBalance).Scan(ctx, &userInfo)

	return &userInfo, err
}
