package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/RaucousDave/streakly/ent"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var EntDB *ent.Client

func DbConnect() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}

	var err error
	connectionString := os.Getenv("NEON_DATABASE_URL")
	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Set to Info level for development
	})

	if err != nil {
		log.Panic(err)
	}

	sqlDb, err := DB.DB()

	if err != nil {
		log.Fatalf("Something went wrong %v", err)
	}

	// Verify connection
	if err := sqlDb.Ping(); err != nil {
		log.Fatalf("Error trying to ping the db %v", err)
	}

	EntDB, err = ent.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Error trying to connect via ent: %v", err)
	}

	if err := EntDB.Schema.Create(context.Background()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected successfully")
}
