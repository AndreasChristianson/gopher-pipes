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

// Just returns a [Source] from the provided items.
func Just[T any](data ...T) Source[T] {
	return FromSlice(data)
}

// FromSlice returns a [Source] from the provided slice of items.
func FromSlice[T any](data []T) Source[T] {
	ret := literalSource[T]{
		data: data,
	}
	ret.log(Verbose, "Creating Source from items(%.10s).", data)
	ret.setStart(ret.start)
	return &ret
}
