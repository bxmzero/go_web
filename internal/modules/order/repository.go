package order

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

func (r *Repository) List(ctx context.Context, userID *int64) ([]Order, error) {
	query := `SELECT id, user_id, item, amount FROM orders`
	args := []any{}
	if userID != nil {
		query += ` WHERE user_id = ?`
		args = append(args, *userID)
	}
	query += ` ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query orders: %w", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var item Order
		if err := rows.Scan(&item.ID, &item.UserID, &item.Item, &item.Amount); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		orders = append(orders, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orders: %w", err)
	}

	return orders, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Order, error) {
	var item Order
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, item, amount FROM orders WHERE id = ?`,
		id,
	).Scan(&item.ID, &item.UserID, &item.Item, &item.Amount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get order by id: %w", err)
	}

	return &item, nil
}

func (r *Repository) Create(ctx context.Context, input CreateOrderRequest) (*Order, error) {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO orders (user_id, item, amount) VALUES (?, ?, ?)`,
		input.UserID,
		input.Item,
		input.Amount,
	)
	if err != nil {
		return nil, fmt.Errorf("insert order: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get inserted order id: %w", err)
	}

	return r.GetByID(ctx, id)
}
