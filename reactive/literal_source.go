package reactive

import (
	"errors"
	"github.com/AndreasChristianson/gopher-pipes/reactive/base-source"
)

type literalSource[T any] struct {
	base_source.BaseSource[T]
	data []T
}

func (l *literalSource[T]) start() {
	defer l.Complete()
	for _, item := range l.data {
		l.Pump(item)
	}
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
	ret.SetStart(ret.start)
	return &ret
}
