package reactive

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
func Buffer[T any](source Source[T], length int) Source[T] {
	c := make(chan T, length)
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
