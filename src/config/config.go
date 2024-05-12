package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
)

func Load() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	Host = os.Getenv("DB_HOST")
	Username = os.Getenv("DB_USER")
	Password = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")

	Port, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		Port = 5432
	}
}
