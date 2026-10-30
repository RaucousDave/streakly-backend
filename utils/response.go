package utils

import (
	"github.com/gofiber/fiber/v3"
)

type Envelope struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorBody `json:"error,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Details any    `json:"details,omitempty"`
}

func Success(c fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(Envelope{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Failure(c fiber.Ctx, status int, code string, message string, details any) error {
	return c.Status(status).JSON(Envelope{
		Success: false,
		Message: message,
		Error:   &ErrorBody{Code: code, Details: details},
	})
}

func OK(c fiber.Ctx, data any) error {
	return Success(c, fiber.StatusOK, "", data)
}

func Created(c fiber.Ctx, message string, data any) error {
	return Success(c, fiber.StatusCreated, message, data)
}

func BadRequest(c fiber.Ctx, message string, details any) error {
	return Failure(c, fiber.StatusBadRequest, "BAD_REQUEST", message, details)
}

func Unauthorized(c fiber.Ctx, message string, details any) error {
	return Failure(c, fiber.StatusUnauthorized, "UNAUTHORIZED", message, details)
}

func NotFound(c fiber.Ctx, message string, details any) error {
	return Failure(c, fiber.StatusNotFound, "NOT_FOUND", message, details)
}

func Conflict(c fiber.Ctx, message string) error {
	return Failure(c, fiber.StatusConflict, "CONFLICT", message, nil)
}

func Internal(c fiber.Ctx, err error) error {
	return Failure(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Something went wrong", err)
}
