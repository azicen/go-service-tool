package transaction

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// SqlxTransactionKey 事务协调器中Sqlx事务唯一键
type SqlxTransactionKey struct{}

type SqlxTransaction struct {
	tx *sqlx.Tx
}

// Commit Sqlx事务提交处理函数
func (tx *SqlxTransaction) Commit(ctx context.Context) error {
	orch, ok := ctx.Value(OrchestratorContextKey{}).(*Orchestrator)
	if !ok {
		return nil
	}

	g, ok := orch.tx[SqlxTransactionKey{}].(*SqlxTransaction)
	if !ok {
		return nil
	}
	return g.tx.Commit()
}

// Rollback Sqlx事务回滚处理函数
func (tx *SqlxTransaction) Rollback(ctx context.Context) error {
	orch, ok := ctx.Value(OrchestratorContextKey{}).(*Orchestrator)
	if !ok {
		return nil
	}

	g, ok := orch.tx[SqlxTransactionKey{}].(*SqlxTransaction)
	if !ok {
		return nil
	}
	return g.tx.Rollback()
}

func SqlxTx(ctx context.Context, db *sqlx.DB) (*sqlx.Tx, error) {
	// 从上下文中获取事务协调器
	orch, ok := ctx.Value(OrchestratorContextKey{}).(*Orchestrator)
	if !ok {
		return db.Beginx()
	}

	// 从事务协调器中获取Sqlx事务
	g, ok := orch.tx[SqlxTransactionKey{}].(*SqlxTransaction)
	if !ok {
		// 不存在则开始新的Sqlx事务
		tx, err := db.Beginx()
		if err != nil {
			return tx, err
		}
		g = &SqlxTransaction{
			tx: tx,
		}
	}
	return g.tx, nil
}

func AddSqlxTransaction(db *sqlx.DB) OrchestratorOption {
	return func(o *Orchestrator) {
		tx, err := db.Beginx()
		if err != nil {
			return
		}
		o.tx[SqlxTransactionKey{}] = &SqlxTransaction{
			tx: tx,
		}
	}
}
