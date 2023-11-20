package reactive

// Map observes one [Source], transform the items observed with the provided mapper function,
// and returns a [Source] of the transformed items. If the mapper returns an error the item dropped, it is not retried.
func Map[T any, V any](source Source[T], mapper func(T) (V, error)) Source[V] {
	c := make(chan V)
	ret := fromChan(c)
	source.UponClose(func() {
		ret.log(Debug, "Closing mapping chan (%p)", c)
		close(c)
		ret.AwaitCompletion()
	})
	source.Observe(func(item T) error {
		defer ret.logPanic(mapper)
		transformed, err := mapper(item)
		if err != nil {
			ret.log(Warning, "Error mapping item (%.10s): [%v]", item, err)
			return err
		}
		ret.log(Verbose, "Mapped item (%.10s) to (%p)", item, transformed)
		c <- transformed
		return nil
	})
	ret.log(Debug, "Created mapped source wit mapper (%p).", mapper)
	ret.Start()
	return ret
}

// Buffer observes one [Source], and returns a [Source] backed by a circular buffer with the requested size.
// This is implemented via a channel.
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
