package reactive

import (
	"sync"
)

// Source is a producer of items. A Source can be based on a generator function ([FromGenerator],
// [FromGeneratorWithExponentialBackoff]), from a channel ([FromChan])
// or from a literal ([Just], [FromSlice]).
type Source[T any] interface {
	UponClose(func())
	Observe(Sink[T])
	Cancel() error
}

type baseSource[T any] struct {
	sinks     []func(T) error
	uponClose []func()
	startFunc func()
}

func (b *baseSource[T]) UponClose(hook func()) {
	b.uponClose = append(b.uponClose, sync.OnceFunc(hook))
}
func (b *baseSource[T]) SetStart(start func()) {
	if b.startFunc != nil {
		panic("startFunc already set")
	}
	b.startFunc = sync.OnceFunc(start)
}

func (b *baseSource[T]) Observe(sink Sink[T]) {
	b.sinks = append(b.sinks, sink)
	if len(b.sinks) == 1 {
		go b.startFunc()
	}
}

func (b *baseSource[T]) complete() {
	for _, hook := range b.uponClose {
		hook()
	}
}
func (b *baseSource[T]) pump(item T) {
	for _, sink := range b.sinks {
		err := sink(item)
		if err != nil {
			//log.Print("potential retry point: ", err)
			//todo
		}
	}
}
