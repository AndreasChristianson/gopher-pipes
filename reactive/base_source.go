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
	logger(Debug, "Registering shutdown hook", hook)
	b.uponClose = append(b.uponClose, sync.OnceFunc(hook))
}
func (b *baseSource[T]) setStart(start func()) {
	logger(Debug, "Registering startup hook", start)
	b.startFunc = sync.OnceFunc(start)
}

func (b *baseSource[T]) Start() {
	logger(Info, "Starting source")
	go func() {
		defer b.complete()
		b.startFunc()
	}()
}

func (b *baseSource[T]) Observe(sink func(T) error) {
	logger(Debug, "Registering sink", sink)
	b.sinks = append(b.sinks, sink)
}

func (b *baseSource[T]) complete() {
	logger(Debug, "Source closed. Cleaning up..")
	for _, hook := range b.uponClose {
		logger(Verbose, "Processing shutdown hook", hook)
		hook()
	}
	logger(Debug, "Cleanup complete.")
}
func (b *baseSource[T]) String() string {
	return fmt.Sprintf("Source=%p", b)
}

func (b *baseSource[T]) pump(item T) {
	logger(Verbose, "Beginning to send item.", item)
	for _, sink := range b.sinks {
		sendItem(item, sink)
	}
	logger(Verbose, "Finished sending item")
}
func sendItem[T any](item T, sink func(T) error) {
	logger(Verbose, "Sending to sink", sink)
	err := sink(item)
	if err != nil {
		logger(Warning, "Failed to write item to sink.", item, sink, err)
	}
}
