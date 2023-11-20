package reactive

import (
	"fmt"
	"sync"
)

type baseSource[T any] struct {
	sinks          []func(T) error
	uponClose      []func()
	startFunc      func()
	lock           sync.Mutex
	closing        bool
	completionLock sync.Mutex
}

func (b *baseSource[T]) UponClose(hook func()) {
	hookOnce := sync.OnceFunc(hook)
	b.lock.Lock()
	b.log(Debug, "Registering shutdown hook (%p) as %d", hook, len(b.uponClose))
	defer b.lock.Unlock()
	b.uponClose = append(b.uponClose, hookOnce)
}
func (b *baseSource[T]) setStart(start func()) {
	b.startFunc = sync.OnceFunc(start)
}

func (b *baseSource[T]) log(level Level, formatString string, args ...interface{}) {
	logger(level, b.String(), formatString, args...)
}

func (b *baseSource[T]) Start() {
	b.log(Info, "Starting source.")
	go func() {
		defer b.complete()
		b.startFunc()
	}()
	b.log(Verbose, "Blocking AwaitCompletion callers.")
	b.completionLock.Lock()
}

func (b *baseSource[T]) Observe(sink func(T) error) {
	b.log(Debug, "Registering sink (%p)", sink)
	b.lock.Lock()
	defer b.lock.Unlock()
	b.sinks = append(b.sinks, sink)
}

func (b *baseSource[T]) complete() {
	b.log(Verbose, "Marking Source as closed.")
	b.closing = true
	b.log(Verbose, "Running %d UponClose hooks..", len(b.uponClose))
	wg := sync.WaitGroup{}
	for index, hook := range b.uponClose {
		wg.Add(1)
		go func(hookToRun func(), indexToRun int) {
			defer wg.Done()
			b.log(Verbose, "Processing UponClose hook %d", indexToRun)
			defer b.logPanic(hookToRun)
			hookToRun()
		}(hook, index)
	}
	wg.Wait()
	b.log(Info, "Source is closed.")
	b.log(Verbose, "Unblocking AwaitCompletion callers.")
	b.completionLock.Unlock()
}
func (b *baseSource[T]) AwaitCompletion() {
	b.completionLock.Lock()
	b.completionLock.Unlock()
}

func (b *baseSource[T]) String() string {
	return fmt.Sprintf("%p", b)
}

func (b *baseSource[T]) pump(item T) {
	if b.closing {
		b.log(Warning, "Ignoring item (%.10s). This source is closing.", item)
		return
	}
	wg := sync.WaitGroup{}
	b.log(Verbose, "Beginning to send item (%.10s)", item)
	for _, sink := range b.sinks {
		wg.Add(1)
		go func(sinkToSendTo func(T) error) {
			defer wg.Done()
			b.sendItem(item, sinkToSendTo)
		}(sink)
	}
	wg.Wait()
	b.log(Verbose, "Finished sending item (%.10s)", item)
}

func (b *baseSource[T]) logPanic(risk interface{}) {
	err := recover()
	if err != nil {
		b.log(Error, "Panic from (%p)! [%v]", risk, err)
	}

}

func (b *baseSource[T]) sendItem(item T, sink func(T) error) {
	b.log(Verbose, "Sending item (%.10s) to sink (%p)", item, sink)
	defer b.logPanic(sink)
	err := sink(item)
	if err != nil {
		b.log(Warning, "Failed to write item (%.10s) to sink (%p): [%s]", item, sink, err)
	}
}
