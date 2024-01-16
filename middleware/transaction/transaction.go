package transaction

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
)

// HandlerFunc 事务中间件的处理函数
type HandlerFunc func(ctx context.Context)

// Option 事务中间件的Option
type Option func() (commit HandlerFunc, rollback HandlerFunc)

// TransactionMap 事务集用于存放各种事务状态
type TransactionMap map[struct{}]any

// TransactionContextKey 上下文的事务唯一键
type TransactionContextKey struct{}

// Transaction 用于处理事务提交和回滚的中间件
// Option 采用先进先出的形式
func Transaction(opts ...Option) middleware.Middleware {
	commitList := make([]HandlerFunc, len(opts))
	rollbackList := make([]HandlerFunc, len(opts))
	for _, o := range opts {
		commit, rollback := o()
		commitList = append(commitList, commit)
		rollbackList = append(rollbackList, rollback)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// 初始化事务
			ctx = context.WithValue(ctx, TransactionContextKey{}, make(TransactionMap))

			// 传递处理
			reply, err = handler(ctx, req)

			if err == nil {
				// 提交事务
				for _, f := range commitList {
					f(ctx)
				}
			} else {
				// 回滚事务
				for _, f := range rollbackList {
					f(ctx)
				}
			}
			return
		}
	}
}
