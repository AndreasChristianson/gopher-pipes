package reactive

// Map observes one [Source], transform the items observed with the provided mapper function,
// and returns a [Source] of the transformed items. The returned [Source] is active immediately.
// Note that this source cannot be cancelled [Source.Cancel]. It closes when the upstream [Source] closes.
func Map[T any, V any](source Source[T], mapper func(T) (V, error)) Source[V] {
	c := make(chan V)
	ret := FromChan(c)
	source.UponClose(func() {
		logger(Debug, "Closing mapping chan.", c)
		close(c)
	})
	source.Observe(func(item T) error {
		transformed, err := mapper(item)
		if err != nil {
			logger(Warning, "Error mapping item.", item, err)
			return err
		}
		logger(Verbose, "Mapped item.", item, transformed)
		c <- transformed
		return nil
	})
	logger(Debug, "Starting mapper.")
	ret.Start()
	return ret
}

// Buffer observes one [Source], and returns a [Source] backed by a circular buffer with the requested size.
// Implemented via a channel.
// Note that this source cannot be cancelled [Source.Cancel]. It closes when the upstream [Source] closes.
func Buffer[T any](source Source[T], size int) Source[T] {
	c := make(chan T, size)
	ret := FromChan(c)
	source.UponClose(func() {
		logger(Debug, "Closing buffered chan.", c)
		close(c)
	})
	source.Observe(func(item T) error {
		c <- item
		return nil
	})
	logger(Debug, "Starting buffer.")
	ret.Start()
	return ret
}
