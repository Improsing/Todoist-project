package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Improsing/go-final-project/db"
	"github.com/Improsing/go-final-project/handlers"
)

var webDir = "./web/"

func main() {
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("/api/task", handlers.TaskHandler)
	http.HandleFunc("/api/tasks", handlers.TasksListHandler)
	http.HandleFunc("/api/task/done", handlers.TaskDoneHandler)


	dbFilePath := db.GetDBFilePath()
	db.CreateDatabase(dbFilePath) 

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "8080"
	}

	if _, err := strconv.Atoi(port); err != nil {
		log.Fatal(err)
	}

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}