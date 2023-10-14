package reactive

type Source[T any] interface {
	Close()
	UponClose(func()) Source[T]
	Observe(Sink[T]) Source[T]
}
type source[T any] struct {
	c chan T
	//observed  bool
	sinks     []func(T)
	uponClose []func()
}

func Just[T any](data ...T) Source[T] {
	return FromSlice(data)
}

func FromSlice[T any](data []T) Source[T] {
	ret := new[T]()
	go func() {
		for _, item := range data {
			ret.c <- item
		}
		ret.Close()
	}()
	return ret
}

func fromChan[T any](c chan T) *source[T] {
	ret := source[T]{
		c:         c,
		uponClose: make([]func(), 0),
		sinks:     make([]func(T), 0),
	}
	go func() {
		for item := range ret.c {
			for _, sink := range ret.sinks {
				sink(item)
			}
		}
	}()
	return &ret
}

func FromChan[T any](c chan T) Source[T] {
	return fromChan[T](c)
}

func new[T any]() *source[T] {
	c := make(chan T)
	return fromChan(c)
}

func FromGenerator[T any](f func() (*T, error)) Source[T] {
	ret := new[T]()
	closed := false
	go func() {
		for {
			if closed {
				return
			}
			item, err := f()
			if err != nil {
				// no more items
				ret.Close()
				return
			}
			if item == nil {
				continue
			}
			select {
			case ret.c <- *item:
			default: // chan closed
				return
			}
		}
	}()
	ret.UponClose(func() {
		closed = true
	})

	return ret
}
