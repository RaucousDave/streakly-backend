package utils

import (
	"errors"
	"time"

	"github.com/RaucousDave/streakly/config"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 30 * 24 * time.Hour

	TypeAccess  = "access"
	TypeRefresh = "refresh"
)

type Claims struct {
	UserId    string `json:"user_id"`
	TokenType string `json:"type"`
	jwt.RegisteredClaims
}

func secret() []byte {
	return []byte(config.Cfg.JWTSecret)
}

// func GenerateToken(id string) (string, error) {
// 	if err := godotenv.Load(); err != nil {
// 		log.Println("Error loading environment variables ", err)
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"user_id": id,
// 	})

// 	t, err := token.SignedString(os.Getenv("JWT_SECRET"))
// 	if err != nil {
// 		return "", err
// 	}
// 	return t, nil
// }

// func VerifyToken(tokenString string) (bool, error) {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
// 		return []byte(os.Getenv("JWT_SECRET")), nil
// 	})
// 	if err != nil {
// 		return false, err
// 	}

// 	return token.Valid, nil
// }

func generateToken(userId, tokenType string, ttl time.Duration) (string, error) {
	claims := Claims{
		UserId:    userId,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret())
}

func GenerateTokenPair(userId string) (access, refresh string, err error) {
	if access, err = generateToken(userId, TypeAccess, AccessTokenTTL); err != nil {
		return "", "", err
	}
	if refresh, err = generateToken(userId, TypeRefresh, RefreshTokenTTL); err != nil {
		return "", "", err
	}

	return
}

func VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return secret(), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
