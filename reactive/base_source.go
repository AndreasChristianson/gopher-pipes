package reactive

import (
	"fmt"
	"sync"
)

type baseSource[T any] struct {
	sinks     []func(T) error
	uponClose []func()
	startFunc func()
}

func (b *baseSource[T]) UponClose(hook func()) {
	b.log(Debug, "Registering shutdown hook", hook)
	b.uponClose = append(b.uponClose, sync.OnceFunc(hook))
}
func (b *baseSource[T]) setStart(start func()) {
	b.log(Debug, "Registering startup hook", start)
	b.startFunc = sync.OnceFunc(start)
}

func (b *baseSource[T]) log(level Level, args ...interface{}) {
	logger(level, b.String(), args...)
}

func (b *baseSource[T]) Start() {
	b.log(Info, "Starting source")
	go func() {
		defer b.complete()
		b.startFunc()
	}()
}

func (b *baseSource[T]) Observe(sink func(T) error) {
	b.log(Debug, "Registering sink", b, sink)
	b.sinks = append(b.sinks, sink)
}

func (b *baseSource[T]) complete() {
	b.log(Debug, "Source closed. Cleaning up..")
	for _, hook := range b.uponClose {
		b.log(Verbose, "Processing shutdown hook", hook)
		hook()
	}
	b.log(Debug, "Cleanup complete.")
}
func (b *baseSource[T]) String() string {
	return fmt.Sprintf("%p", b)
}

func (b *baseSource[T]) pump(item T) {
	b.log(Verbose, "Beginning to send item:", item)
	for _, sink := range b.sinks {
		b.sendItem(item, sink)
	}
	b.log(Verbose, "Finished sending item")
}

func (b *baseSource[T]) sendItem(item T, sink func(T) error) {
	b.log(Verbose, "Sending to sink", sink)
	err := sink(item)
	if err != nil {
		b.log(Warning, "Failed to write item to sink.", item, sink, err)
	}
}
