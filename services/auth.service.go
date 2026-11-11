package services

import (
	"context"
	"errors"

	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/ent"
	"github.com/RaucousDave/streakly/ent/user"
)

type AuthService struct {
	db *ent.Client
}

func NewAuthService(db *ent.Client) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) UpsertGoogleUser(ctx context.Context, googleId string, email, name, picture, timezone string) (*ent.User, error) {
	u, err := s.db.User.Query().Where(user.GoogleID(googleId)).Only(ctx)

	if err == nil {
		return u, nil
	}

	if !ent.IsNotFound(err) {
		return nil, err
	}

	u, err = s.db.User.Query().Where(user.Email(email)).Only(ctx)

	if err == nil {
		return u, nil
	}

	if !ent.IsNotFound(err) {
		return nil, err
	}

	u, err = s.db.User.Create().SetGoogleID(googleId).SetEmail(email).SetDisplayName(name).SetAvatarURL(picture).SetAuthProvider("GOOGLE").Save(ctx)

	if ent.IsConstraintError(err) {
		return s.db.User.Query().Where(user.Or(user.GoogleID(googleId), user.Email(email))).Only(ctx)
	}

	return u, err
}

func (s *AuthService) Login(c context.Context, dto dto.LoginDto) (*ent.User, error) {
	u, err := s.db.User.Query().Where(user.Email(dto.Email)).Only(c)

	if ent.IsNotFound(err) {
		return nil, errors.New("User does not exist")
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *AuthService) Register(c context.Context, dto dto.AuthRequest) (*ent.User, error) {
	u, err := s.db.User.Create().SetDisplayName(dto.DisplayName).SetEmail(dto.Email).SetPasswordHash(dto.Password).SetAuthProvider("EMAIL").SetTimezone(dto.Timezone).Save(c)

	if ent.IsConstraintError(err) {
		return nil, err
	}

	return u, err
}
