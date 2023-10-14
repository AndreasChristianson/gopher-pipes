package reactive

func (s *source[T]) Observe(sink Sink[T]) Source[T]{
	s.sinks = append(s.sinks, sink)
	return s
}
