# Task API

REST API для управления задачами.

## Стек

- Go
- net/http
- PostgreSQL
- pgx / pgxpool
- JSON
- Postman
- pgAdmin

## Возможности

- Создать задачу
- Получить список задач
- Получить задачу по id
- Обновить задачу
- Удалить задачу

## Архитектура

handler -> service -> repository -> PostgreSQL

## Endpoints

GET /ping  
GET /tasks  
POST /tasks  
GET /tasks/{id}  
PATCH /tasks/{id}  
DELETE /tasks/{id}

## Переменные окружения

DB_URL — строка подключения к PostgreSQL.

Пример:

```bash
DB_URL='postgres://user:password@localhost:5432/task_api?sslmode=disable'

Запуск:

DB_URL='postgres://user:password@localhost:5432/task_api?sslmode=disable' go run ./cmd/app
``` 

Тесты:

go test ./...

БД:

migrations/001_create_tasks.sql