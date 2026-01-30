package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppID     string
	AppSecret string
	Port      string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	appID := os.Getenv("LARK_APP_ID")
	appSecret := os.Getenv("LARK_APP_SECRET")
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	if appID == "" || appSecret == "" {
		log.Fatal("LARK_APP_ID and LARK_APP_SECRET must be set")
	}

	return &Config{
		AppID:     appID,
		AppSecret: appSecret,
		Port:      port,
	}
}
