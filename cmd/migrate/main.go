package main

import (
	"locations-project/internal/app/ds"
	"locations-project/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(
		&ds.User{},
		&ds.Location{},
		&ds.PlayersLocationGame{},
		&ds.PlayersChosenLocation{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
