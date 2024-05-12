package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Host      string
	Port      int
	Username  string
	Password  string
	DBName    string
	SecretKey []byte
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
	SecretKey = []byte(os.Getenv("SECRET_KEY"))

	Port, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		Port = 5432
	}
}
