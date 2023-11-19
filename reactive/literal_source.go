package reactive

type literalSource[T any] struct {
	baseSource[T]
	data []T
}

func (l *literalSource[T]) start() {
	for _, item := range l.data {
		l.pump(item)
	}
}

// Just returns a [Source] from the provided items. The returned [Source] is active immediately.
func Just[T any](data ...T) Source[T] {
	return FromSlice(data)
}

// FromSlice returns a [Source] from the provided slice of items. The returned [Source] is active immediately.
func FromSlice[T any](data []T) Source[T] {
	logger(Verbose, "Creating Source from items.", data)
	ret := literalSource[T]{
		data: data,
	}
	ret.setStart(ret.start)
	return &ret
}
