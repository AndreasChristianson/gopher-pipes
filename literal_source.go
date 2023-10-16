package reactive

import "errors"

type literalSource[T any] struct {
	baseSource[T]
	data []T
}

func (l *literalSource[T]) start() {
	go func() {
		defer l.complete()
		for _, item := range l.data {
			l.pump(item)
		}
	}()
}
func (l *literalSource[T]) Cancel() error {
	return errors.New(
		"this source is based on a literal values. " +
			"it will clode when all values are emitted",
	)
}

func Just[T any](data ...T) Source[T] {
	return FromSlice(data)
}

func FromSlice[T any](data []T) Source[T] {
	ret := literalSource[T]{
		data:     data,
	}
	ret.start()
	return &ret
}
