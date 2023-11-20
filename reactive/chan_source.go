package reactive

type chanSource[T any] struct {
	channel chan T
	baseSource[T]
}

func (c *chanSource[T]) start() {
	for item := range c.channel {
		c.pump(item)
	}
}

// FromChan returns a [Source] from the provided channel. The returned [Source] is active immediately.
func FromChan[T any](channel chan T) Source[T] {
	return fromChan(channel)
}

func fromChan[T any](channel chan T) *chanSource[T] {
	ret := chanSource[T]{
		channel: channel,
	}
	ret.log(Verbose, "Creating chan based Source: chan(%p)", channel)
	ret.setStart(ret.start)
	return &ret
}
