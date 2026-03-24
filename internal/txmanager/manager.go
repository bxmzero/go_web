package txmanager

import (
	"context"

	"gorm.io/gorm"
)

type txContextKey struct{}

var txKey = txContextKey{}

type Manager interface {
	WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
	DB(ctx context.Context) *gorm.DB
	Current(ctx context.Context) (*gorm.DB, bool)
}

type manager struct {
	db *gorm.DB
}

func New(db *gorm.DB) Manager {
	return &manager{db: db}
}

func withTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func extractTx(ctx context.Context) (*gorm.DB, bool) {
	if ctx == nil {
		return nil, false
	}
	tx, ok := ctx.Value(txKey).(*gorm.DB)
	return tx, ok
}

func (m *manager) WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if tx, ok := extractTx(ctx); ok && tx != nil {
		return fn(ctx)
	}

	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(withTx(ctx, tx))
	})
}

func (m *manager) DB(ctx context.Context) *gorm.DB {
	if tx, ok := extractTx(ctx); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	if ctx == nil {
		return m.db
	}
	return m.db.WithContext(ctx)
}

func (m *manager) Current(ctx context.Context) (*gorm.DB, bool) {
	return extractTx(ctx)
}
