package reactive

import (
	"github.com/AndreasChristianson/gopher-pipes/reactive/base-source"
	"math"
	"time"
)

type generatorSource[T any] struct {
	generator func() *GeneratorResponse[T]
	base_source.BaseSource[T]
	cancelled             bool
	maxBackoff            int
	minBackoff            int
	consecutiveErrorCount int
}

func (g *generatorSource[T]) start() {
	defer g.Complete()
	for {
		if g.cancelled {
			return
		}
		response := g.generator()

		if response.Data != nil {
			g.Pump(*response.Data)
		}
		if response.Finished {
			return
		}
		g.incrementError(response.Err)
		g.exponentialBackoff()
	}
}

func (g *generatorSource[T]) Cancel() error {
	g.cancelled = true
	return nil
}

func (g *generatorSource[T]) exponentialBackoff() {
	if g.consecutiveErrorCount == 0 && g.maxBackoff == 0 {
		return
	}
	backOff := float64(g.minBackoff) * math.Pow(2.0, float64(g.consecutiveErrorCount))
	millisecondsToWait := time.Duration(
		min(float64(g.maxBackoff), backOff),
	)
	<-time.After(millisecondsToWait * time.Millisecond)
}

func (g *generatorSource[T]) incrementError(err error) {
	if err != nil {
		g.consecutiveErrorCount++
	} else {
		g.consecutiveErrorCount = 0
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
func FromGenerator[T any](generator func() *GeneratorResponse[T]) Source[T] {
	return FromGeneratorWithExponentialBackoff(generator, 0, 0)
}

// FromGeneratorWithExponentialBackoff is similar to FromGenerator, but accepts parameters for implementing a exponential
// backoff to prevent rapid polling.
// - maxBackoff maximum time to wait in milliseconds
// - minBackoff minimum time to wait in milliseconds
func FromGeneratorWithExponentialBackoff[T any](generator func() *GeneratorResponse[T], maxBackoff int, minBackoff int) Source[T] {
	ret := generatorSource[T]{
		generator:  generator,
		maxBackoff: maxBackoff,
		minBackoff: minBackoff,
	}
	ret.SetStart(ret.start)
	return &ret
}
