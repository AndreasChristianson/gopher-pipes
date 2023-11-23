package reactive

import (
	"errors"
	"math"
	"time"
)

type generatorSource[T any] struct {
	generator func() (*T, error)
	baseSource[T]
	maxBackoff            float64
	backoffMultiplier     float64
	consecutiveErrorCount int
}

func (g *generatorSource[T]) start() {
	for !g.closing {
		g.log(Verbose, "Polling generator (%p).", g.generator)
		response, err := g.generator()
		if response != nil {
			g.pump(*response)
		}
		if err != nil {
			var t *GeneratorFinished
			switch {
			case errors.As(err, &t):
				g.log(Debug, "Generator(%p) func indicates completion.", g.generator)
				return
			default:
				g.consecutiveErrorCount++
				g.log(Info, "Error from generator: [%v]", err)
				g.log(Debug, "Error count incremented: %d", g.consecutiveErrorCount)
				g.exponentialBackoff()
			}
			continue
		}
		g.clearErrorCount()
	}
}

func (g *generatorSource[T]) Cancel() error {
	g.log(Info, "Cancel request received. Marking source as closed.")
	g.closing = true
	return nil
}

func (g *generatorSource[T]) exponentialBackoff() {
	if g.consecutiveErrorCount == 0 || g.maxBackoff == 0 {
		return
	}
	backOff := g.backoffMultiplier * math.Pow(2.0, float64(g.consecutiveErrorCount))
	wait := time.Duration(min(g.maxBackoff, backOff)) * time.Millisecond

	g.log(Verbose, "Waiting %s before next generator poll.", wait)
	time.Sleep(wait)
}

func (g *generatorSource[T]) clearErrorCount() {
	if g.consecutiveErrorCount > 0 {
		g.log(Debug, "Clearing error count")
		g.consecutiveErrorCount = 0
	}
}

// FromGenerator returns a [Source] from the provided generator function. The generator provided will be polled
// for items:
//   - The generator should return (*T, nil) when an item is available.
//   - If no item is available, the generator should return (nil, nil).
//   - If the generator is complete it should return (*T, GeneratorFinished)
//   - Finally, errors may be returned via (nil, error).
//
// Note that this source can be cancelled via [CancellableSource.Cancel].
func FromGenerator[T any](generator func() (*T, error)) CancellableSource[T] {
	return FromGeneratorWithExponentialBackoff(generator, 0, 0)
}

// FromGeneratorWithDefaultBackoff is similar to FromGenerator, but waits at least 250ms and at most 10s
func FromGeneratorWithDefaultBackoff[T any](generator func() (*T, error)) CancellableSource[T] {
	return FromGeneratorWithExponentialBackoff(generator, 10000, 125)
}

// FromGeneratorWithExponentialBackoff is similar to FromGenerator, but accepts parameters for implementing a exponential
// backoff to prevent rapid polling.
//   - maxBackoff maximum time to wait in milliseconds
//   - backoffMultiplier the multiplier m in m*2^e where e is the error count
func FromGeneratorWithExponentialBackoff[T any](
	generator func() (*T, error),
	maxBackoff float64,
	backoffMultiplier float64,
) CancellableSource[T] {
	ret := generatorSource[T]{
		generator:         generator,
		maxBackoff:        maxBackoff,
		backoffMultiplier: backoffMultiplier,
	}
	ret.log(
		Verbose,
		"Creating Source with exp backoff: Generator (%p), maxBackoff (%.0fms), backoffMultiplier (%.0fms)",
		generator,
		maxBackoff,
		backoffMultiplier,
	)
	ret.setStart(ret.start)
	return &ret
}
