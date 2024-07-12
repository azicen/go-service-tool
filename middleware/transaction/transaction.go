package transaction

import (
	"context"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
)

// Transaction 事务
type Transaction interface {
	Commit(context.Context) error
	Rollback(context.Context) error
}

// OrchestratorContextKey 上下文的事务协调器唯一键
type OrchestratorContextKey struct{}

// Orchestrator 事务协调器
type Orchestrator struct {
	Tx *orderedmap.OrderedMap[struct{}, Transaction] // 有序的 map[struct{}]Transaction

	log *log.Helper
}

func newOrchestrator() *Orchestrator {
	return &Orchestrator{
		Tx: orderedmap.NewOrderedMap[struct{}, Transaction](),
	}
}

// OrchestratorOption 事务协调器选项设置
type OrchestratorOption func(*Orchestrator)

// Logger 设置事务协调器日志
func Logger(logger log.Logger) OrchestratorOption {
	return func(o *Orchestrator) {
		o.log = log.NewHelper(logger)
	}
}

// AddTransaction 初始化时立刻添加事务
func AddTransaction(key struct{}, tx Transaction) OrchestratorOption {
	return func(o *Orchestrator) {
		o.Tx.Set(key, tx)
	}
}

// Middleware 用于处理事务提交和回滚的中间件
// Option 采用先进先出的形式
func Middleware(opts ...OrchestratorOption) middleware.Middleware {
	orch := newOrchestrator()
	for _, o := range opts {
		o(orch)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// 事务管理器添加到上下文
			ctx = context.WithValue(ctx, OrchestratorContextKey{}, orch)

			// 传递处理
			reply, err = handler(ctx, req)

			if err == nil {
				// 提交事务
				for el := orch.Tx.Front(); el != nil; el = el.Next() {
					err := el.Value.Commit(ctx)
					orch.log.Error(err)
				}
			} else {
				// 回滚事务
				for el := orch.Tx.Front(); el != nil; el = el.Next() {
					err := el.Value.Rollback(ctx)
					orch.log.Error(err)
				}
			}
			return
		}
	}
}
