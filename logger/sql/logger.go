package sql

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	sqldblogger "github.com/simukti/sqldb-logger"
)

// SqlLogger SQL 日志接口的实现
type SqlLogger struct {
	*log.Helper
}

// NewLogger 创建一个新的 LoggerHelper 实例
func NewLogger(l log.Logger) *SqlLogger {
	h := log.NewHelper(l)
	return &SqlLogger{h}
}

func (s *SqlLogger) Log(_ context.Context, level sqldblogger.Level, msg string, data map[string]interface{}) {
	var lvl log.Level

	switch level {
	case sqldblogger.LevelError:
		lvl = log.LevelError
	case sqldblogger.LevelInfo:
		lvl = log.LevelInfo
	case sqldblogger.LevelDebug:
		lvl = log.LevelDebug
	case sqldblogger.LevelTrace:
		lvl = log.LevelDebug
	default:
		lvl = log.LevelDebug
	}

	keyvals := make([]interface{}, 0, len(data)*2+2)
	keyvals = append(keyvals, log.DefaultMessageKey, msg)
	for k, v := range data {
		keyvals = append(keyvals, k, v)
	}

	s.Helper.Log(lvl, keyvals...)
}
