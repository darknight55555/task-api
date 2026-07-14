package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"task-api/internal/db"
	"task-api/internal/handler"
	"task-api/internal/repository"
	"task-api/internal/service"
)

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DB_URL")
	if strings.TrimSpace(dbURL) == "" {
		log.Fatal("DB_URL is required")
	}

	pool, errConnect := db.ConnectPostgres(ctx, dbURL)
	if errConnect != nil {
		log.Fatal(errConnect)
	}

	log.Println("database connected")

	defer pool.Close()

	repo := repository.NewPostgresTaskRepository(pool)
	serv := service.NewTaskService(repo)
	taskHandler := handler.NewTaskHandler(serv)

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", taskHandler.HandlePing)
	mux.HandleFunc("/tasks", taskHandler.HandleTasks)
	mux.HandleFunc("/tasks/", taskHandler.HandleTaskByID)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
