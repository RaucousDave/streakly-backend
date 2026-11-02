package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	JWTSecret          string
	NeonDatabaseUrl    string
	GoogleClientId     string
	MinAppVersion      string
	LatestAppVersion   string
	MaintenanceMode    string
	UpdateUrl          string
	MaintenanceMessage string
}

var Cfg Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Print("Error loading environment variables", err)
	}

	Cfg = Config{
		Port:               getEnv("PORT", "5000"),
		NeonDatabaseUrl:    mustGetEnv("NEON_DATABASE_URL"),
		GoogleClientId:     mustGetEnv("GOOGLE_CLIENT_ID"),
		JWTSecret:          mustGetEnv("JWT_SECRET"),
		MinAppVersion:      getEnv("MIN_APP_VERSION", "1.0.0"),
		LatestAppVersion:   getEnv("LATEST_APP_VERSION", "1.0.0"),
		MaintenanceMode:    getEnv("MAINTENANCE_MODE", "false"),
		UpdateUrl:          getEnv("UPDATE_URL", "https://streakly-nu.vercel.app/"),
		MaintenanceMessage: getEnv("MAINTENANCE_MESSAGE", "We're doing some upgrades. Back shortly"),
	}
}

func getEnv(key string, fallback string) string {
	v := os.Getenv(key)

	if v == "" {
		return fallback
	}

	return v
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)

	if v == "" {
		log.Fatalf("Error loading env variables for %v", key)
	}

	return v
}
