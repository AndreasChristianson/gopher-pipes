package reactive

func (s *source[T]) UponClose(f func()) Source[T] {
	s.uponClose = append(s.uponClose, f)
	return s
}
