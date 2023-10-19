package reactive

// Map observes one [Source], transform the items observed with the provided mapper function,
// and returns a [Source] of the transformed items. The returned [Source] is active immediately.
// Note that this source cannot be cancelled [Source.Cancel]. It closes when the upstream [Source] closes.
func Map[T any, V any](source Source[T], mapper func(T) (V, error)) Source[V] {
	c := make(chan V)
	ret := fromChan(c)
	source.UponClose(func() {
		close(c)
	})
	source.Observe(func(item T) error {
		transformed, err := mapper(item)
		if err != nil {
			return err
		}
		c <- transformed
		return nil
	})
	return ret
}

// Buffer observes one [Source], and returns a [Source] backed by a circular buffer with the requested size.
// Implemented via a channel.
// Note that this source cannot be cancelled [Source.Cancel]. It closes when the upstream [Source] closes.
func Buffer[T any](source Source[T], size int) Source[T] {
	c := make(chan T, size)
	ret := fromChan(c)
	source.UponClose(func() {
		close(c)
	})
	source.Observe(func(item T) error {
		c <- item
		return nil
	})
	return ret
}
