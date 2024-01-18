package unittest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

type TestLogger struct {
	*testing.T
}

func (t TestLogger) Log(level log.Level, keyvals ...interface{}) error {
	var builder strings.Builder
	builder.WriteString(level.String())

	// 解析keyvals奇数元素为key，偶数元素为value，形成key=value的形式加入builder
	for i := 0; i < len(keyvals); i += 2 {
		key := keyvals[i]
		value := keyvals[i+1]
		builder.WriteString(fmt.Sprintf(" %v=%v", key, value))
	}

	t.Logf(builder.String())
	return nil
}

func MokeLogger(t *testing.T) log.Logger {
	return TestLogger{t}
}
