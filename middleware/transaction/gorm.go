package transaction

import (
	"context"

	"gorm.io/gorm"
)

// GORMTransactionContextKey 上下文的GORM事务唯一键
type GORMTransactionContextKey struct{}

func GORMTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	// 从上下文中获取事务数据集
	m, ok := ctx.Value(TransactionContextKey{}).(TransactionMap)
	if !ok {
		return db.Begin()
	}

	// 从事务数据集中获取gorm事务
	tx, ok := m[GORMTransactionContextKey{}].(*gorm.DB)
	if !ok {
		// 不存在则开始新的GORM事务
		tx = db.Begin()
	}
	return tx
}

// GORMCommitHandler Gorm事务提交处理函数
func GORMCommitHandler(ctx context.Context) {
	m, ok := ctx.Value(TransactionContextKey{}).(TransactionMap)
	if !ok {
		return
	}

	tx, ok := m[GORMTransactionContextKey{}].(*gorm.DB)
	if !ok {
		return
	}
	tx.Commit()
}

// GORMRollbackHandler Gorm事务回滚处理函数
func GORMRollbackHandler(ctx context.Context) {
	m, ok := ctx.Value(TransactionContextKey{}).(TransactionMap)
	if !ok {
		return
	}

	tx, ok := m[GORMTransactionContextKey{}].(*gorm.DB)
	if !ok {
		return
	}

	tx.Rollback()
}

func GORMTransactionOption() Option {
	return func() (commit HandlerFunc, rollback HandlerFunc) {
		return GORMCommitHandler, GORMRollbackHandler
	}
}
