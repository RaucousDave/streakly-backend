package handlers

import (
	"context"
	"time"

	"github.com/RaucousDave/streakly/config"
	"github.com/RaucousDave/streakly/database"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/mod/semver"
)

type AppConfigResponse struct {
	Status              string `json:"status"`
	ServerTime          string `json:"server_time"`
	MinSupportedVersion string `json:"min_supported_version"`
	LatestVersion       string `json:"latest_version"`
	ForceUpdate         bool   `json:"force_update"`
	UpdateAvailable     bool   `json:"update_available"`
	UpdateUrl           string `json:"update_url"`
	Maintenance         bool   `json:"maintenance"`
	MaintenanceMessage  string `json:"maintenance_message,omitempty"`
}

func GetAppConfig(c fiber.Ctx) error {
	clientVersion := c.Query("version")
	minVersion := config.Cfg.MinAppVersion
	latestVersion := config.Cfg.LatestAppVersion

	status := "ok"
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	if sqlDb, err := database.DB.DB(); err != nil || sqlDb.PingContext(ctx) != nil {
		status = "degraded"
	}

	res := AppConfigResponse{
		Status:              status,
		ServerTime:          time.Now().UTC().Format(time.RFC3339),
		MinSupportedVersion: minVersion,
		LatestVersion:       latestVersion,
		UpdateUrl:           config.Cfg.UpdateUrl,
		Maintenance:         config.Cfg.MaintenanceMode == "true",
		MaintenanceMessage:  config.Cfg.MaintenanceMessage,
	}

	if semver.IsValid("v" + clientVersion) {
		res.ForceUpdate = semver.Compare("v"+clientVersion, "v"+minVersion) < 0
		res.UpdateAvailable = semver.Compare("v"+clientVersion, "v"+latestVersion) < 0
	}

	return c.JSON(res)
}
