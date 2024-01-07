package grom

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/logger"
)

// GormLoggerHelper GORM Logger 接口的实现
type GormLoggerHelper struct {
	*log.Helper
}

// NewLoggerHelper 创建一个新的 NacosLoggerHelper 实例
func NewLoggerHelper(l log.Logger) *GormLoggerHelper {
	h := log.NewHelper(l)
	return &GormLoggerHelper{h}
}

// LogMode 设置日志级别
func (g *GormLoggerHelper) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

// Info 日志信息级别
func (g *GormLoggerHelper) Info(ctx context.Context, msg string, data ...interface{}) {
	g.Infof(msg, data...)
}

// Warn 日志警告级别
func (g *GormLoggerHelper) Warn(ctx context.Context, msg string, data ...interface{}) {
	g.Warnf(msg, data...)
}

// Error 日志错误级别
func (g *GormLoggerHelper) Error(ctx context.Context, msg string, data ...interface{}) {
	g.Errorf(msg, data...)
}

// Trace 用于跟踪 SQL 执行和记录慢查询
func (g *GormLoggerHelper) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		g.Errorw("err", err, "elapsed", elapsed, "sql", sql, "rows", rows)
	} else {
		g.Debugw("elapsed", elapsed, "sql", sql, "rows", rows)
	}
}
