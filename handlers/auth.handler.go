package handlers

import (
	"errors"
	"time"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/RaucousDave/streakly/config"
	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/services"
	"github.com/RaucousDave/streakly/utils"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req dto.AuthRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequest(c, "invalid request body", nil)
	}

	passwordHash := utils.HashPassword(req.Password)

	user := dto.AuthRequest{
		Email:       req.Email,
		Password:    passwordHash,
		DisplayName: req.DisplayName,
		Timezone:    req.Timezone,
	}

	newUser, err := h.authService.Register(c, user)
	if err != nil {
		return utils.Internal(c, err)
	}

	return utils.Created(c, "User created successfully", newUser)
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req dto.LoginDto

	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body", nil)
	}

	user, err := h.authService.Login(c, req)

	if err != nil {
		return utils.BadRequest(c, "Invalid credentials", nil)
	}

	if user == nil || user.PasswordHash == nil {
		return utils.Internal(c, errors.New("Login returned a nil user"))
	}

	if !utils.VerifyPassword(*user.PasswordHash, req.Password) {
		return utils.BadRequest(c, "Invalid credentials", nil)
	}

	access, refresh, err := utils.GenerateTokenPair(user.ID.String())

	if err != nil {
		return utils.Internal(c, err)
	}

	token := dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}

	return utils.OK(c, token)
}

func (h *AuthHandler) GoogleLogin(c fiber.Ctx) error {
	var body dto.GoogleLoginDto
	

	if err := c.Bind().Body(&body); err != nil || body.IDToken == "" {
		return utils.BadRequest(c, "idToken required", nil)
	}

	payload, err := idtoken.Validate(c.Context(), body.IDToken, config.Cfg.GoogleClientId)
	if err != nil {
		return utils.Unauthorized(c, "Unauthorized", nil)
	}

	if verified, _ := payload.Claims["email_verified"].(bool); !verified {
		return utils.Unauthorized(c, "Unauthorized", nil)
	}

	timezone := body.Timezone
	if body.Timezone == ""{
			timezone = "UTC"
	}

	
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	
	if _, err := time.LoadLocation(timezone); err != nil{
		return utils.BadRequest(c, "Invalid timezone", nil)
	}

	user, err := h.authService.UpsertGoogleUser(c, payload.Subject, email, name, picture, timezone)
	if err != nil {
		return utils.Internal(c, err)
	}

	access, refresh, err := utils.GenerateTokenPair(user.ID.String())
	if err != nil {
		return utils.Internal(c, err)
	}

	return utils.OK(c, dto.AuthResponse{AccessToken: access, RefreshToken: refresh})
}
