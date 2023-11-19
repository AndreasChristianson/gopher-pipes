package reactive

// Source is a producer of items. A Source can be based on a generator function ([FromGenerator],
// [FromGeneratorWithExponentialBackoff]), from a channel ([FromChan])
// or from a literal ([Just]).
type Source[T any] interface {
	UponClose(func())
	Observe(func(T) error)
	Start()
}

type CancellableSource[T any] interface {
	Source[T]
	Cancel() error
}
