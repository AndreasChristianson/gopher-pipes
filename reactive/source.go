package reactive

// Source is a producer of items. A Source can be based on a generator function ([FromGenerator],
// [FromGeneratorWithExponentialBackoff]), from a channel ([FromChan])
// or from a literal ([Just]).
type Source[T any] interface {
	// UponClose registers a hook to run then this source shuts down.
	// All registered functions will complete before AwaitCompletion unblocks
	UponClose(func())
	// Observe registers a sink that will observe each item that flows through this source.
	// Each sink will be called in its own go routine.
	Observe(func(T) error)
	// Start begins pumping items through the source.
	// Generators start polling, channels start listening, literals start pumping.
	Start()
	// AwaitCompletion blocks until the source is closed and all UponClose hooks are complete.
	AwaitCompletion()
}

type CancellableSource[T any] interface {
	Source[T]
	Cancel() error
}
