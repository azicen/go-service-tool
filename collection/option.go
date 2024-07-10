// 来自 https://github.com/noxiouz/golang-generics-util/blob/main/collection/option.go

package collection

// Option 代表一个可选的值
type Option[T any] interface {
	// HasValue 如果`Option`包含一个实际的值，则`HasValue`属性为`true`
	HasValue() bool
	// Value 返回存储的值
	// 如果`HasValue`为`false`，则结果未定义
	// 它可能会返回零初始化的`T`类型值
	Value() T
}

// None 没有值的`Option`
func None[T any]() Option[T] {
	return noneImpl[T]{}
}

type noneImpl[T any] struct{}

func (noneImpl[T]) HasValue() bool {
	return false
}

func (noneImpl[T]) Value() T {
	v := new(T)
	return *v
}

// Some 有值的`Option`
func Some[T any](value T) Option[T] {
	return &someImpl[T]{
		value: value,
	}
}

type someImpl[T any] struct {
	value T
}

func (s *someImpl[T]) HasValue() bool {
	return true
}

func (s *someImpl[T]) Value() T {
	return s.value
}

// Map applies a function to an optional value
func Map[T any, U any](opt Option[T], fn func(T) U) Option[U] {
	if !opt.HasValue() {
		return None[U]()
	}

	return Some(fn(opt.Value()))
}

// FlatMap applies a function to an optional value. Similar to Map, but the function
// returns an Option
func FlatMap[T any, U any](opt Option[T], fn func(T) Option[U]) Option[U] {
	if !opt.HasValue() {
		return None[U]()
	}

	return fn(opt.Value())
}
