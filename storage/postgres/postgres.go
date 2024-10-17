package postgres

import (
	"fmt"
	"os"
	"test-case/internal/config"
	"test-case/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Database *gorm.DB
}

func New(config config.Config) (*Database, error) {
	const op = "storage.postgres.New"

	var dbUrl string

	if config.Env == "local" {
		dbUrl = fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			config.Host, config.Port, config.User, config.Password, config.DbName)
	} else {
		dbUrl = os.Getenv("DATABASE_URL")
	}

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if config.Env != "prod" {
		db = db.Debug()
	}

	if err := db.AutoMigrate(&models.Group{}, &models.Song{}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Database{Database: db}, nil
}

func (db *Database) Stop() error {
	const op = "storage.postgres.Stop"

	storage, err := db.Database.DB()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	storage.Close()

	return nil
}
