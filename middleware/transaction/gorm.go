package transaction

import (
	"context"

	"gorm.io/gorm"
)

// GORMTransactionKey 事务协调器中GORM事务唯一键
type GORMTransactionKey struct{}

type GORMTransaction struct {
	tx *gorm.DB
}

// Commit Gorm事务提交处理函数
func (tx *GORMTransaction) Commit(ctx context.Context) error {
	orch, ok := ctx.Value(OrchestratorContextKey{}).(*Orchestrator)
	if !ok {
		return nil
	}

	g, ok := orch.tx[GORMTransactionKey{}].(*GORMTransaction)
	if !ok {
		return nil
	}
	return g.tx.Commit().Error
}

// Rollback GORM事务回滚处理函数
func (tx *GORMTransaction) Rollback(ctx context.Context) error {
	orch, ok := ctx.Value(OrchestratorContextKey{}).(*Orchestrator)
	if !ok {
		return nil
	}

	g, ok := orch.tx[GORMTransactionKey{}].(*GORMTransaction)
	if !ok {
		return nil
	}
	return g.tx.Rollback().Error
}

func GORMTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	// 从上下文中获取事务协调器
	orch, ok := ctx.Value(OrchestratorContextKey{}).(*Orchestrator)
	if !ok {
		return db.Begin()
	}

	// 从事务协调器中获取gorm事务
	g, ok := orch.tx[GORMTransactionKey{}].(*GORMTransaction)
	if !ok {
		// 不存在则开始新的GORM事务
		g = &GORMTransaction{
			tx: db.Begin(),
		}
	}
	return g.tx
}

func AddGORMTransaction(db *gorm.DB) OrchestratorOption {
	return func(o *Orchestrator) {
		o.tx[GORMTransactionKey{}] = &GORMTransaction{
			tx: db.Begin(),
		}
	}
}
