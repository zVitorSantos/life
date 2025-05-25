package config

import (
	"fmt"
	"os"

	"life/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB inicializa a conexão com o banco de dados
func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migra as tabelas
	err = db.AutoMigrate(&models.User{}, &models.APIKey{}, &models.RefreshToken{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
