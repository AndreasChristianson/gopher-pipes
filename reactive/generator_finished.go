package reactive

// GeneratorFinished should be returned from a generator function when the generator completes successfully,
// but no further items are expected.
type GeneratorFinished struct{} //todo: should generators just return three things?

// Error implements the error interface
func (g *GeneratorFinished) Error() string {
	return "Generator finished"
}
