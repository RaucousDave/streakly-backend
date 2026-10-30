package services

import (
	"context"
	"errors"

	"github.com/RaucousDave/streakly/database"
	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/ent"
	"github.com/RaucousDave/streakly/ent/user"
)

func UpsertGoogleUser(ctx context.Context, googleId string, email, name, picture string) (*ent.User, error) {
	u, err := database.EntDB.User.Query().Where(user.GoogleID(googleId)).Only(ctx)

	if err == nil {
		return u, nil
	}

	if !ent.IsNotFound(err) {
		return nil, err
	}

	u, err = database.EntDB.User.Query().Where(user.Email(email)).Only(ctx)

	if err == nil {
		return u, nil
	}

	if !ent.IsNotFound(err) {
		return nil, err
	}

	u, err = database.EntDB.User.Create().SetGoogleID(googleId).SetEmail(email).SetDisplayName(name).SetAvatarURL(picture).SetAuthProvider("GOOGLE").Save(ctx)

	if ent.IsConstraintError(err) {
		return database.EntDB.User.Query().Where(user.Or(user.GoogleID(googleId), user.Email(email))).Only(ctx)
	}

	return u, err
}

func Login(c context.Context, dto dto.LoginDto) (*ent.User, error) {
	u, err := database.EntDB.User.Query().Where(user.Email(dto.Email)).Only(c)

	if ent.IsNotFound(err) {
		return nil, errors.New("User does not exist")
	}

	return u, nil
}

func Register(c context.Context, dto dto.AuthRequest) (*ent.User, error) {
	u, err := database.EntDB.User.Create().SetDisplayName(dto.DisplayName).SetEmail(dto.Email).SetPasswordHash(dto.Password).SetAuthProvider("EMAIL").SetTimezone(dto.Timezone).Save(c)

	if ent.IsConstraintError(err) {
		return nil, err
	}

	return u, err
}
