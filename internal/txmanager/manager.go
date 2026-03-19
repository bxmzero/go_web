package txmanager

import (
	"context"

	"gorm.io/gorm"
)

type contextKey string

const txKey contextKey = "txmanager:current"

type Manager interface {
	WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
	DB(ctx context.Context) *gorm.DB
}

type manager struct {
	db *gorm.DB
}

func New(db *gorm.DB) Manager {
	return &manager{db: db}
}

func (m *manager) WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey, tx)
		return fn(txCtx)
	})
}

func (m *manager) DB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx.WithContext(ctx)
	}

	return m.db.WithContext(ctx)
}
