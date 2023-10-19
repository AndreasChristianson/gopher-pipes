package reactive

import "errors"

type chanSource[T any] struct {
	c chan T
	baseSource[T]
}

func (c *chanSource[T]) Cancel() error {
	return errors.New(
		"this source is based on a chanel." +
			" it will clode when the channel closes",
	)
}

func fromChan[T any](c chan T) *chanSource[T] {
	ret := chanSource[T]{
		c: c,
	}
	go func() {
		defer ret.complete()
		for item := range ret.c {
			ret.pump(item)
		}
	}()
	return &ret

}

// FromChan returns a [Source] from the provided channel. The returned [Source] is active immediately.
// Note that this source cannot be cancelled [Source.Cancel]. It closes when the channel closes.
func FromChan[T any](c chan T) Source[T] {
	return fromChan[T](c)
}
