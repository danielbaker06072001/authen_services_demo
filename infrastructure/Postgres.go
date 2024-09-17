package infrastructure

import (
	"authen-service/appConfig/config"

	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect connects to a PostgreSQL database using GORM.
func Connect(cfg *config.AppConfig) (*gorm.DB, *gorm.DB, error) {
	dsn := cfg.Postgres.DB_URL_Authen
	dsn1 := cfg.Postgres.DB_URL_DATA

	db_authen, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, nil, err
	}

	db_data, err := gorm.Open(postgres.Open(dsn1), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, nil, err
	}

	log.Println("connection established")
	return db_authen, db_data, nil
}
