package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
}

func NewConfig() *Config {
	loadEnv()
	return &Config{
		ServerAddress: fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
	}
}

func loadEnv() {
	for _, p := range []string{".env", "../../.env"} {
		if err := godotenv.Load(p); err == nil {
			return
		}
	}
	log.Println("No .env file found")
}