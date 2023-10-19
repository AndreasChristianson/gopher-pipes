package reactive

type Sink[T any] func(T) error

//type sink[T any] struct {
//	lock         sync.Mutex
//	registerFunc func(*baseSource[T])
//	observeFunc  func(T)
//	registered   bool
//}
//
//func (s *sink[T]) register(source *baseSource[T]) error {
//	s.lock.Lock()
//	defer s.lock.Unlock()
//	if s.registered{
//		return errors.New("sink already registered")
//	}
//	s.registerFunc(source)
//	s.registered = true
//	return nil
//}
//
//func (s *sink[T]) observe(item T) {
//	s.observeFunc(item)
//}
//
//func NewSink[T any](observeFunc func(T)) Sink[T] {
//	return &sink[T]{
//		observeFunc: observeFunc,
//		registerFunc: func(s *baseSource[T]) {},
//	}
//}
