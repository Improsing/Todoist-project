package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func GetDBFilePath() string {
	dbFilePath := os.Getenv("TODO-DBFILE")
	if dbFilePath == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		dbFilePath = filepath.Join(currentDir, "scheduler.db")
	}
	return dbFilePath
}


func CreateDatabase(dbFilePath string) {
	// Открываем соединение с БД
	var err error
	DB, err = sql.Open("sqlite", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	//defer DB.Close()

	// Создаем таблицу scheduler
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT CHECK(length(repeat) <= 128)
	);`)
	
	if err != nil {
		log.Fatal(err)
	}

	// Создаем индекс по полю date
	_, err = DB.Exec(`CREATE INDEX IF NOT EXISTS date_idx ON scheduler(date);`)
	if err != nil {
		log.Fatal(err)
	}
}
