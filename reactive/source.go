package reactive

// Source is a producer of items. A Source can be based on a generator function ([FromGenerator],
// [FromGeneratorWithExponentialBackoff]), from a channel ([FromChan])
// or from a literal ([Just], [FromSlice]).
type Source[T any] interface {
	UponClose(func())
	Observe(func(T) error)
	Cancel() error
	// todo, consider Start()
}