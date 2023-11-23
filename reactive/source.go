package reactive

// Source is a producer of items. A Source can be based on a generator function ([FromGenerator],
// [FromGeneratorWithExponentialBackoff]), a channel ([FromChan])
// or a literal ([Just]).
type Source[T any] interface {
	// UponClose registers a hook to run then this source shuts down.
	// All registered functions will complete before AwaitCompletion unblocks.
	UponClose(func())
	// Observe registers a sink that will observe each item that flows through this source.
	// Each sink will be called in its own go routine.
	Observe(func(T) error)
	// Start begins pumping items through the source.
	// Generators start polling, channels start listening, literals start pumping.
	//
	// This method can be called multiple times, only the first has any effect (locking via sync.OnceFunc).
	Start()
	// AwaitCompletion blocks until the source is closed and all UponClose hooks are complete.
	AwaitCompletion()
}

// CancellableSource is a [Source] that can be canceled.
// or a literal ([Just]).
type CancellableSource[T any] interface {
	Source[T]
	// Cancel stops the propagation of items from this source. If the [CancellableSource] is based on a generator, the
	// generator is no longer polled.
	Cancel() error
}
