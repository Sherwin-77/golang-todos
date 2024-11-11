package main

import (
	"log"

	"github.com/sherwin-77/go-echo-template/configs"
	"github.com/sherwin-77/go-echo-template/internal/entity"
	"github.com/sherwin-77/go-echo-template/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	config := configs.LoadConfig()
	db, err := database.InitDB(config.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	users := []entity.User{
		{
			Username: "admin",
			Email:    "admin@example.com",
			Password: string(password),
		},
		{
			Username: "editor",
			Email:    "editor@example.com",
			Password: string(password),
		},
		{
			Username: "user",
			Email:    "user@example.com",
			Password: string(password),
		},
	}

	db.Create(&users)

	log.Println("Users table seeded successfully")
}
