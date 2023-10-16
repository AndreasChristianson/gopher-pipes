package reactive

import (
	"log"
	"math"
	"time"
)

type generatorSource[T any] struct {
	generator func() *GeneratorResponse[T]
	baseSource[T]
	cancelled             bool
	maxExponentialBackoff int
	exponentialFactor     float64
	consecutiveErrorCount int
}

func (g *generatorSource[T]) start() {
	go func() {
		defer g.complete()
		for {
			if g.cancelled {
				return
			}
			response := g.generator()

			if response.Data != nil {
				g.pump(*response.Data)
			}
			if response.Finished {
				return
			}
			g.incrementError(response.Err)
			g.exponentialBackoff()
		}
	}()
}

func (g *generatorSource[T]) Cancel() error {
	g.cancelled = true
	return nil
}

func (g *generatorSource[T]) exponentialBackoff() {
	backOff := min(g.maxExponentialBackoff, g.consecutiveErrorCount)
	if backOff == 0 {
		return
	}

	millisecondsToWait := time.Duration(
		max(math.Pow(g.exponentialFactor, float64(backOff)), 100),
	)
	<-time.After(millisecondsToWait * time.Millisecond)

}

func (g *generatorSource[T]) incrementError(err error) {
	if err != nil {
		g.consecutiveErrorCount++
		if g.consecutiveErrorCount == g.maxExponentialBackoff {
			log.Print("warning, reached maximum exponential backoff: ", err)
		}
	} else {
		if g.maxExponentialBackoff > 0 && g.consecutiveErrorCount >= g.maxExponentialBackoff {
			log.Print("exponential backoff warning resolved")
		}
		g.consecutiveErrorCount = 0
	}
}

func FromGenerator[T any](generator func() *GeneratorResponse[T]) Source[T] {
	return FromGeneratorWithExponentialBackoff(generator, 0, 0)
}
func FromGeneratorWithExponentialBackoff[T any](generator func() *GeneratorResponse[T], maxExponentialBackoff int, exponentialFactor float64) Source[T] {
	ret := generatorSource[T]{
		generator:             generator,
		maxExponentialBackoff: maxExponentialBackoff,
		exponentialFactor:     exponentialFactor,
	}
	ret.start()
	return &ret
}
