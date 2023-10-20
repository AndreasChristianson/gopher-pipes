package reactive

import (
	"errors"
	"github.com/AndreasChristianson/gopher-pipes/reactive/base-source"
)

type chanSource[T any] struct {
	channel chan T
	base_source.BaseSource[T]
}

func (c *chanSource[T]) Cancel() error {
	return errors.New(
		"this source is based on a chanel." +
			" it will clode when the channel closes",
	)
}

func (c *chanSource[T]) start() {
	defer c.Complete()
	for item := range c.channel {
		c.Pump(item)
	}
}

func fromChan[T any](channel chan T) *chanSource[T] {
	ret := chanSource[T]{
		channel: channel,
	}
	ret.SetStart(ret.start)
	return &ret

}

// FromChan returns a [Source] from the provided channel. The returned [Source] is active immediately.
// Note that this source cannot be cancelled [Source.Cancel]. It closes when the channel closes.
func FromChan[T any](c chan T) Source[T] {
	return fromChan[T](c)
}
