package repository

import (
	"context"
	"errors"
	"task-api/internal/model"

	"github.com/jackc/pgx/v5"
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
// Create — QueryRow + INSERT RETURNING
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
// List — Query + rows.Next
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
// GetByID — QueryRow + SELECT WHERE
func (p *PostgresTaskRepository) GetByID(ctx context.Context, id int) (model.Task, error) {
	var task model.Task

	sqlQuery := `
		SELECT id, title, done, created_at 
		FROM tasks
		WHERE id = $1;
	`

	err := p.pool.QueryRow(ctx, sqlQuery, id).Scan(
		&task.ID,
		&task.Title,
		&task.Done,
		&task.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Task{}, ErrTaskNotFound
		}

		return model.Task{}, err
	}

	return task, nil
}
// Update — QueryRow + UPDATE RETURNING
func (p *PostgresTaskRepository) Update(
	ctx context.Context,
	id int,
	title string,
	done bool,
) (model.Task, error) {
	var task model.Task

	sqlQuery := `
			UPDATE tasks 
			SET title = $1, done = $2
			WHERE id = $3
			RETURNING id, title, done, created_at;
		`

	err := p.pool.QueryRow(ctx, sqlQuery, title, done, id).Scan(
		&task.ID,
		&task.Title,
		&task.Done,
		&task.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Task{}, ErrTaskNotFound
		}

		return model.Task{}, err
	}

	return task, nil
}
// Delete — QueryRow + DELETE RETURNING id
func (p *PostgresTaskRepository) Delete(
	ctx context.Context,
	id int,
) error {
	var deletedID int

	sqlQuery := `
			DELETE FROM tasks
			WHERE id = $1
			RETURNING id;
		`

	err := p.pool.QueryRow(ctx, sqlQuery, id).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}

		return err
	}

	return nil
}
