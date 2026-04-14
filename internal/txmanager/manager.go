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
//ctx：外部传入的上下文（context.Context）。
//fn：一个回调函数，类型是 func(txCtx context.Context) error，它在事务中执行实际的数据库操作。
//回调函数 fn 是传递给 WithinTransaction 的实际操作（例如数据库查询、更新等）,在调用的时候需要实现的匿名函数
//这个函数的作用就是，让匿名函数fn的入参txCtx一定是携带了事务DB的上下文
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
