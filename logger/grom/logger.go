package grom

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/logger"
)

// GormLogger GORM 日志接口的实现
type GormLogger struct {
	*log.Helper
}

// NewLogger 创建一个新的日志实例
func NewLogger(l log.Logger) *GormLogger {
	h := log.NewHelper(l)
	return &GormLogger{h}
}

// LogMode 设置日志级别
func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

// Info 日志信息级别
func (g *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	g.Infof(msg, data...)
}

// Warn 日志警告级别
func (g *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	g.Warnf(msg, data...)
}

// Error 日志错误级别
func (g *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	g.Errorf(msg, data...)
}

// Trace 用于跟踪 SQL 执行和记录慢查询
func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		g.Errorw("err", err, "elapsed", elapsed, "sql", sql, "rows", rows)
	} else {
		g.Debugw("elapsed", elapsed, "sql", sql, "rows", rows)
	}
}
