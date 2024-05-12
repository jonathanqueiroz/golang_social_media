// database/database.go
package db

import (
	"database/sql"
	"fmt"

	"project01/src/config"

	_ "github.com/lib/pq"
)

func New() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	fmt.Println("Conexão bem sucedida!")
	return db, nil
}
