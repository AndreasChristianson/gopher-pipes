package reactive

func Map[T any, V any](source Source[T], mapper func(T) V) Source[V] {
	c := make(chan V)
	ret := FromChan(c)
	source.Observe(func(item T) {
		transformed := mapper(item)
		select {
		case c <- transformed:
		default: // chan closed
			return
		}
	})
	return ret
}
