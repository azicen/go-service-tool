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

// Orchestrator 协调器
type Orchestrator struct {
	// Tx 有序的 map[struct{}]Transaction
	Tx *orderedmap.OrderedMap[struct{}, Transaction]
}

func newOrchestrator() *Orchestrator {
	return &Orchestrator{
		Tx: orderedmap.NewOrderedMap[struct{}, Transaction](),
	}
}

// GetOrchestratorFromContext 从上下文中获取事务协调器
func GetOrchestratorFromContext(ctx context.Context) *Orchestrator {
	orch, ok := ctx.Value(OrchestratorContextKey{}).(*Orchestrator)
	if !ok {
		// 不应该出现的情况
		panic("")
	}
	return orch
}

// manager 事务协调器
type manager struct {
	log         *log.Helper
	initHandler []InitHandler
}

func newManager() *manager {
	return &manager{
		initHandler: make([]InitHandler, 0, 10),
	}
}

// ManagerOption 事务管理器选项设置
type ManagerOption func(*manager)

// Logger 设置事务协调器日志
func Logger(logger log.Logger) ManagerOption {
	return func(m *manager) {
		m.log = log.NewHelper(logger)
	}
}

// 初始化事务处理器函数
type InitHandler func(*Orchestrator)

// AddInitHandler 初始化时立刻添加事务
func AddInitHandler(f InitHandler) ManagerOption {
	return func(m *manager) {
		m.initHandler = append(m.initHandler, f)
	}
}

// Middleware 用于处理事务提交和回滚的中间件。
// opts: ManagerOption 采用先进先出的形式
func Middleware(opts ...ManagerOption) middleware.Middleware {
	m := newManager()
	for _, o := range opts {
		o(m)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			orch := newOrchestrator()
			for _, h := range m.initHandler {
				h(orch)
			}
			// 事务管理器添加到上下文
			ctx = context.WithValue(ctx, OrchestratorContextKey{}, orch)

			// 传递处理
			return handler(ctx, req)
		}
	}
}

// Commit 提交事务
func Commit(ctx context.Context) error {
	orch := GetOrchestratorFromContext(ctx)
	for el := orch.Tx.Front(); el != nil; el = el.Next() {
		err := el.Value.Commit(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Rollback 回滚事务
func Rollback(ctx context.Context) error {
	orch := GetOrchestratorFromContext(ctx)
	for el := orch.Tx.Front(); el != nil; el = el.Next() {
		err := el.Value.Rollback(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
