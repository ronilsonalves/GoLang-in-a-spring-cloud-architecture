package config

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var URLDbConnection = ""

// LoadConfig - load info about db connection from env variables
func LoadConfig() {
	var err error
	if err = godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	URLDbConnection = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("DATABASE_URL"),
		os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME"))
}

// ConnectDatabase - perform a sql.Open() to connect to database
func ConnectDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", URLDbConnection)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
