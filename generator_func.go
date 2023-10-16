package reactive

type GeneratorResponse[T any] struct {
	Data *T
	Finished bool
	Err error
}

