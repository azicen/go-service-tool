package strconv

import (
	"fmt"
	"strconv"
)

type Int interface {
	int | int8 | int16 | int32 | int64
}

type UInt interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type Float interface {
	float32 | float64
}

type Number interface {
	Int | UInt | Float
}

// FormatNumber 将数字转换为字符串
func FormatNumber[T Number](i T) (str string, err error) {
	switch v := any(i).(type) {
	case int:
		str = strconv.FormatInt(int64(v), 10)
	case int8:
		str = strconv.FormatInt(int64(v), 10)
	case int16:
		str = strconv.FormatInt(int64(v), 10)
	case int32:
		str = strconv.FormatInt(int64(v), 10)
	case int64:
		str = strconv.FormatInt(v, 10)

	case uint:
		str = strconv.FormatUint(uint64(v), 10)
	case uint8:
		str = strconv.FormatUint(uint64(v), 10)
	case uint16:
		str = strconv.FormatUint(uint64(v), 10)
	case uint32:
		str = strconv.FormatUint(uint64(v), 10)
	case uint64:
		str = strconv.FormatUint(v, 10)

	case float32:
		str = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		str = strconv.FormatFloat(v, 'f', -1, 64)

	default:
		return "", fmt.Errorf("无法转换不支持的类型数据: %v", i)
	}
	return str, nil
}

// ParseNumber 将字符串转换为数字
func ParseNumber[T Number](s string, i *T) error {
	switch any(*i).(type) {
	case int:
		v, err := strconv.ParseInt(s, 10, 0)
		*i = T(v)
		return err
	case int8:
		v, err := strconv.ParseInt(s, 10, 8)
		*i = T(v)
		return err
	case int16:
		v, err := strconv.ParseInt(s, 10, 16)
		*i = T(v)
		return err
	case int32:
		v, err := strconv.ParseInt(s, 10, 32)
		*i = T(v)
		return err
	case int64:
		v, err := strconv.ParseInt(s, 10, 64)
		*i = T(v)
		return err

	case uint:
		v, err := strconv.ParseUint(s, 10, 0)
		*i = T(v)
		return err
	case uint8:
		v, err := strconv.ParseUint(s, 10, 8)
		*i = T(v)
		return err
	case uint16:
		v, err := strconv.ParseUint(s, 10, 16)
		*i = T(v)
		return err
	case uint32:
		v, err := strconv.ParseUint(s, 10, 32)
		*i = T(v)
		return err
	case uint64:
		v, err := strconv.ParseUint(s, 10, 64)
		*i = T(v)
		return err

	case float32:
		v, err := strconv.ParseFloat(s, 32)
		*i = T(v)
		return err
	case float64:
		v, err := strconv.ParseFloat(s, 64)
		*i = T(v)
		return err

	default:
		return fmt.Errorf("无法解析不支持的类型数据: %v", *i)
	}
}
