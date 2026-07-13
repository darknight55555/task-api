package main

import (
	"context"
	"log"
	"net/http"

	"task-api/internal/db"
	"task-api/internal/handler"
	"task-api/internal/repository"
	"task-api/internal/service"
)

func main() {
	ctx := context.Background()
	connString := "postgres://darkxam:1@localhost:5432/task_api?sslmode=disable"

	pool, errConnect := db.ConnectPostgres(ctx, connString)
	if errConnect != nil {
		log.Fatal(errConnect)
		return
	}

	log.Println("database connected")

	defer pool.Close()

	//postgresRepo := repository.NewPostgresTaskRepository(pool)

	repo := repository.NewMemoryTaskRepository()
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
