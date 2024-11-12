package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env       string
	Key       string
	JWTSecret string
	Name      string
	Port      string
	Postgres  PostgresConfig
	Redis     RedisConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func GetConfig() *Config {
	config := &Config{
		Env:       os.Getenv("ENV"),
		Key:       os.Getenv("APP_KEY"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		Name:      os.Getenv("APP_NAME"),
		Port:      os.Getenv("APP_PORT"),
		Postgres: PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DB"),
		},
		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		},
	}

	// Fallback to APP_KEY if JWT_SECRET is not set
	if config.JWTSecret == "" {
		config.JWTSecret = config.Key
	}

	return config
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	return GetConfig()
}
