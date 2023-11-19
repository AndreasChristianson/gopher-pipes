package reactive

import (
	"sync"
)

type asyncSource[T any] struct {
	baseSource[T]
	wg sync.WaitGroup
}

func (a *asyncSource[T]) sendToAll(item T) {
	logger(Verbose, "Beginning to send item async.", item)
	for _, sink := range a.sinks {
		a.wg.Add(1)
		go a.send(item, sink)
	}
}

func (a *asyncSource[T]) send(item T, sink func(T) error) {
	defer a.wg.Done()
	sendItem(item, sink)
}

// Async observes one [Source], and pumps items to its observers in asynchronously.
// Observers may encounter items in any order.
func Async[T any](source Source[T]) Source[T] {
	ret := &asyncSource[T]{}
	source.Observe(func(item T) error {
		go ret.sendToAll(item)
		return nil
	})
	source.UponClose(func() {
		logger(Verbose, "Awaiting async cleanup..")
		ret.wg.Wait()
		logger(Verbose, "Async cleanup Complete.")
		ret.complete()
	})
	return ret
}
