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

type GeneratorFinished error

func (g *generatorSource[T]) start() {
	for {
		if g.cancelled {
			logger(Debug, "Discovered that source is cancelled.")
			return
		}
		response, err := g.generator()
		if err != nil {
			var t GeneratorFinished
			switch {
			case errors.As(err, &t):
				logger(Debug, "Generator func indicates completion.")
				return
			default:
				g.consecutiveErrorCount++
				logger(Debug, "Error from generator", err)
				logger(Verbose, "Error count incremented.", g.consecutiveErrorCount)
				g.exponentialBackoff()
			}
		}

		if response != nil {
			g.pump(*response)
		}
	}
}

func (g *generatorSource[T]) Cancel() error {
	logger(Debug, "Marking source as cancelled.")
	g.cancelled = true
	return nil
}

func (g *generatorSource[T]) exponentialBackoff() {
	if g.consecutiveErrorCount == 0 || g.maxBackoff == 0 {
		return
	}
	backOff := g.backoffMultiplier * math.Pow(2.0, float64(g.consecutiveErrorCount))
	wait := time.Duration(min(g.maxBackoff, backOff)) * time.Millisecond

	logger(Verbose, fmt.Sprintf("Waiting %s before next generator poll.", wait))
	time.Sleep(wait)
}

func (g *generatorSource[T]) incrementError(err error) {
	if err != nil {
		g.consecutiveErrorCount++
		logger(Debug, "Error from generator", err)
		logger(Verbose, "Error count incremented.", g.consecutiveErrorCount)
	} else {
		g.consecutiveErrorCount = 0
		logger(Verbose, "Error count reset.")
	}
}

// FromGenerator returns a [Source] from the provided generator function. The
// returned [Source] is active immediately. The generator provided will be polled
// for items, it should return GeneratorResponse{Data: item} when an item is
// available. If no item is available, the generator should return
// GeneratorResponse{Data: nil}. If the generator is complete it should return
// GeneratorResponse{Finished: true}. Finally, errors may be returned via
// GeneratorResponse[]{ Err: err}. Note that this source can be cancelled via
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
	logger(Verbose, "Creating generator based Source with exponential backoff.", generator, maxBackoff, backoffMultiplier)
	ret := generatorSource[T]{
		generator:         generator,
		maxBackoff:        maxBackoff,
		backoffMultiplier: backoffMultiplier,
	}
	ret.setStart(ret.start)
	return &ret
}
