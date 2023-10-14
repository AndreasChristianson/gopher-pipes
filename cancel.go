package reactive

func (s *source[T]) Close() {
	close(s.c)
	for _, hook := range s.uponClose {
		hook()
	}
}
