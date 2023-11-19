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
	logger(Verbose, "Creating chan based Source.", channel)
	ret := chanSource[T]{
		channel: channel,
	}
	ret.setStart(ret.start)
	return &ret
}
