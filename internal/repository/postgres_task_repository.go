package repository

import (
	"context"
	"task-api/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTaskRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresTaskRepository(pool *pgxpool.Pool) *PostgresTaskRepository {
	postgresTaskRepository := PostgresTaskRepository{
		pool: pool,
	}

	return &postgresTaskRepository
}

func (p *PostgresTaskRepository) Create(ctx context.Context, title string) (model.Task, error) {
	var task model.Task

	sqlQuery := `
		INSERT INTO tasks (title)
		VALUES ($1)
		RETURNING id, title, done, created_at;
	`

	err := p.pool.QueryRow(ctx, sqlQuery, title).Scan(
		&task.ID,
		&task.Title,
		&task.Done,
		&task.CreatedAt,
	)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (p *PostgresTaskRepository) List(ctx context.Context) ([]model.Task, error) {
	sqlQuery := `
		SELECT id, title, done, created_at
		FROM tasks
		ORDER BY id;
	`

	rows, errQuery := p.pool.Query(ctx, sqlQuery)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)

	for rows.Next() {
		var task model.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Done,
			&task.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
