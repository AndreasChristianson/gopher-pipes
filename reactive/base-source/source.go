package base_source

import (
	"sync"
)

type BaseSource[T any] struct {
	sinks     []func(T) error
	uponClose []func()
	startFunc func()
}

func (b *BaseSource[T]) UponClose(hook func()) {
	b.uponClose = append(b.uponClose, sync.OnceFunc(hook))
}
func (b *BaseSource[T]) SetStart(start func()) {
	if b.startFunc != nil {
		panic("startFunc already set")
	}
	b.startFunc = sync.OnceFunc(start)
}

func (b *BaseSource[T]) Observe(sink func(T) error) {
	b.sinks = append(b.sinks, sink)
	if len(b.sinks) == 1 {
		go b.startFunc()
	}
}

func (b *BaseSource[T]) Complete() {
	for _, hook := range b.uponClose {
		hook()
	}
}
func (b *BaseSource[T]) Pump(item T) {
	for _, sink := range b.sinks {
		err := sink(item)
		if err != nil {
			//log.Print("potential retry point: ", err)
			//todo
		}
	}
}
