package trans

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// SqlxTransactionKey 事务协调器中Sqlx事务唯一键
type SqlxTransactionKey struct{}

type SqlxTransaction struct {
	dbTx *sqlx.Tx
}

// Commit Sqlx事务提交处理函数
func (tx *SqlxTransaction) Commit(ctx context.Context) error {
	orch := GetOrchestratorFromContext(ctx)

	t, ok := orch.Tx.Get(SqlxTransactionKey{})
	if !ok {
		return nil
	}

	sqlxTx, ok := t.(*SqlxTransaction)
	if !ok {
		return nil
	}
	return sqlxTx.dbTx.Commit()
}

// Rollback Sqlx事务回滚处理函数
func (tx *SqlxTransaction) Rollback(ctx context.Context) error {
	orch := GetOrchestratorFromContext(ctx)

	t, ok := orch.Tx.Get(SqlxTransactionKey{})
	if !ok {
		return nil
	}

	sqlxTx, ok := t.(*SqlxTransaction)
	if !ok {
		return nil
	}
	return sqlxTx.dbTx.Rollback()
}

func SqlxTx(ctx context.Context, db *sqlx.DB) (*sqlx.Tx, error) {
	orch := GetOrchestratorFromContext(ctx)

	// 从事务协调器中获取Sqlx事务
	t, ok := orch.Tx.Get(SqlxTransactionKey{})
	if !ok {
		// 不存在则开始新的Sqlx事务
		tx, err := db.Beginx()
		if err != nil {
			return tx, err
		}
		t = &SqlxTransaction{
			dbTx: tx,
		}
		orch.Tx.Set(SqlxTransactionKey{}, t)
	}

	sqlxTx, _ := t.(*SqlxTransaction)
	return sqlxTx.dbTx, nil
}

func Sqlx(db *sqlx.DB) ManagerOption {
	return AddInitHandler(
		func(o *Orchestrator) {
			tx, err := db.Beginx()
			if err != nil {
				return
			}
			o.Tx.Set(SqlxTransactionKey{}, &SqlxTransaction{
				dbTx: tx,
			})
		},
	)
}
