package reactive

import (
	"sync"
)

type Source[T any] interface {
	UponClose(func())
	Observe(Sink[T])
	Cancel() error
}

type baseSource[T any] struct {
	sinks     []func(T) error
	uponClose []func()
}

func (b *baseSource[T]) UponClose(hook func()) {
	b.uponClose = append(b.uponClose, sync.OnceFunc(hook))
}

func (b *baseSource[T]) Observe(sink Sink[T]) {
	b.sinks = append(b.sinks, sink)
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
