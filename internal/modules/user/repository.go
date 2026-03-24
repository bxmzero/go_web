package user

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, email FROM users ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var item User
		if err := rows.Scan(&item.ID, &item.Name, &item.Email); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*User, error) {
	var item User
	err := r.db.QueryRowContext(ctx, `SELECT id, name, email FROM users WHERE id = ?`, id).
		Scan(&item.ID, &item.Name, &item.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &item, nil
}

func (r *Repository) Create(ctx context.Context, input CreateUserRequest) (*User, error) {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (name, email) VALUES (?, ?)`,
		input.Name,
		input.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get inserted user id: %w", err)
	}

	return r.GetByID(ctx, id)
}
