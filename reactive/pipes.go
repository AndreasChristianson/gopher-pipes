package reactive

// Map observes one [Source], transform the items observed with the provided mapper function,
// and returns a [Source] of the transformed items. The returned [Source] is active immediately.
func Map[T any, V any](source Source[T], mapper func(T) (V, error)) Source[V] {
	c := make(chan V)
	ret := fromChan(c)
	source.UponClose(func() {
		ret.log(Debug, "Closing mapping chan (%p)", c)
		close(c)
		ret.AwaitCompletion()
	})
	source.Observe(func(item T) error {
		transformed, err := mapper(item)
		if err != nil {
			ret.log(Warning, "Error mapping item (%p): [%v]", item, err)
			return err
		}
		ret.log(Verbose, "Mapped item (%p) to (%p)", item, transformed)
		c <- transformed
		return nil
	})
	ret.log(Debug, "Created mapped source.")
	ret.Start()
	return ret
}

// Buffer observes one [Source], and returns a [Source] backed by a circular buffer with the requested size.
// Implemented via a channel.
func Buffer[T any](source Source[T], size int) Source[T] {
	c := make(chan T, size)
	ret := fromChan(c)
	source.UponClose(func() {
		ret.log(Debug, "Closing buffered chan (%p).", c)
		close(c)
		ret.AwaitCompletion()
	})
	source.Observe(func(item T) error {
		c <- item
		return nil
	})
	ret.log(Debug, "Created buffered source.")
	ret.Start()
	return ret
}
