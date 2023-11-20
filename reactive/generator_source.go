package reactive

import (
	"errors"
	"fmt"
	"math"
	"time"
)

type generatorSource[T any] struct {
	generator func() (*T, error)
	baseSource[T]
	cancelled             bool
	maxBackoff            float64
	backoffMultiplier     float64
	consecutiveErrorCount int
}

type GeneratorFinished struct{}

func (g GeneratorFinished) Error() string {
	return "Generator finished"
}

func (g *generatorSource[T]) start() {
	for {
		if g.cancelled {
			g.log(Debug, "Discovered that source is cancelled.")
			return
		}
		response, err := g.generator()
		if response != nil {
			g.pump(*response)
		}
		if err != nil {
			var t GeneratorFinished
			switch {
			case errors.As(err, &t):
				g.log(Debug, "Generator func indicates completion.")
				return
			default:
				g.consecutiveErrorCount++
				g.log(Debug, fmt.Sprintf("Error from generator: [%v]", err))
				g.log(Verbose, "Error count incremented.", g.consecutiveErrorCount)
				g.exponentialBackoff()
			}
			continue
		}
		g.consecutiveErrorCount = 0
	}
}

func (g *generatorSource[T]) Cancel() error {
	g.log(Debug, "Marking source as cancelled.")
	g.cancelled = true
	return nil
}

func (g *generatorSource[T]) exponentialBackoff() {
	if g.consecutiveErrorCount == 0 || g.maxBackoff == 0 {
		return
	}
	backOff := g.backoffMultiplier * math.Pow(2.0, float64(g.consecutiveErrorCount))
	wait := time.Duration(min(g.maxBackoff, backOff)) * time.Millisecond

	g.log(Debug, fmt.Sprintf("Waiting %s before next generator poll.", wait))
	time.Sleep(wait)
}

// FromGenerator returns a [Source] from the provided generator function. The
// returned [Source] is active immediately. The generator provided will be polled
// for items, it should return `*T,nil` when an item is
// available. If no item is available, the generator should return
// `nil,nil`. If the generator is complete it should return
// `*T,GeneratorFinished`. Finally, errors may be returned via
// `nil,err`. Note that this source can be cancelled via
// [Source.Cancel].
func FromGenerator[T any](generator func() (*T, error)) CancellableSource[T] {
	return FromGeneratorWithExponentialBackoff(generator, 0, 0)
}

// FromGeneratorWithDefaultBackoff is similar to FromGenerator, but waits at least 500ms and at most 10s
func FromGeneratorWithDefaultBackoff[T any](generator func() (*T, error)) CancellableSource[T] {
	return FromGeneratorWithExponentialBackoff(generator, 10000, 250)
}

// FromGeneratorWithExponentialBackoff is similar to FromGenerator, but accepts parameters for implementing a exponential
// backoff to prevent rapid polling.
// - maxBackoff maximum time to wait in milliseconds
// - backoffMultiplier the multiplier m in m*2^e where e is the error count
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
	ret.log(Verbose, "Creating generator based Source with exponential backoff.", generator, maxBackoff, backoffMultiplier)
	ret.setStart(ret.start)
	return &ret
}
