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

// Just returns a [Source] from the provided items. The returned [Source] is active immediately.
// Note that this source cannot be cancelled [Source.Cancel]. It closes when all items are observed.
func Just[T any](data ...T) Source[T] {
	return FromSlice(data)
}

// FromSlice returns a [Source] from the provided slice of items. The returned [Source] is active immediately.
// Note that this source cannot be cancelled [Source.Cancel]. It closes when all items are observed.
func FromSlice[T any](data []T) Source[T] {
	ret := literalSource[T]{
		data: data,
	}
	ret.start()
	return &ret
}
