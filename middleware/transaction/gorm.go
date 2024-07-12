package transaction

import (
	"context"

	"gorm.io/gorm"
)

// GORMTransactionKey 事务协调器中GORM事务唯一键
type GORMTransactionKey struct{}

type GORMTransaction struct {
	dbTx *gorm.DB
}

// Commit Gorm事务提交处理函数
func (tx *GORMTransaction) Commit(ctx context.Context) error {
	orch, ok := ctx.Value(OrchestratorContextKey{}).(*Orchestrator)
	if !ok {
		return nil
	}

	t, ok := orch.Tx.Get(GORMTransactionKey{})
	if !ok {
		return nil
	}

	gormTx, ok := t.(*GORMTransaction)
	if !ok {
		return nil
	}
	return gormTx.dbTx.Commit().Error
}

// Rollback GORM事务回滚处理函数
func (tx *GORMTransaction) Rollback(ctx context.Context) error {
	orch, ok := ctx.Value(OrchestratorContextKey{}).(*Orchestrator)
	if !ok {
		return nil
	}

	t, ok := orch.Tx.Get(GORMTransactionKey{})
	if !ok {
		return nil
	}

	gormTx, ok := t.(*GORMTransaction)
	if !ok {
		return nil
	}
	return gormTx.dbTx.Rollback().Error
}

func GORMTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	// 从上下文中获取事务协调器
	orch, ok := ctx.Value(OrchestratorContextKey{}).(*Orchestrator)
	if !ok {
		return db.Begin()
	}

	// 从事务协调器中获取gorm事务
	t, ok := orch.Tx.Get(GORMTransactionKey{})
	if !ok {
		// 不存在则开始新的GORM事务
		t = &GORMTransaction{
			dbTx: db.Begin(),
		}
		orch.Tx.Set(GORMTransactionKey{}, t)
	}
	gormTx, _ := t.(*GORMTransaction)
	return gormTx.dbTx
}

func AddGORMTransaction(db *gorm.DB) OrchestratorOption {
	return func(o *Orchestrator) {
		o.Tx.Set(GORMTransactionKey{}, &GORMTransaction{
			dbTx: db.Begin(),
		})
	}
}
