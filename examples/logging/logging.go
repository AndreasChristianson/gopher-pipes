package main

import (
	"errors"
	"fmt"
	"github.com/AndreasChristianson/gopher-pipes/reactive"
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	reactive.SetLogger(func(level reactive.Level, id string, args ...interface{}) {
		if level < reactive.Debug {
			return
		}
		log.Default().Printf("%s [%s]: %v", level, id, args)
	})
	wg := sync.WaitGroup{}
	returns := []string{"Hello", "world", "!"}
	pos := 0
	wg.Add(1)
	pipe := reactive.FromGeneratorWithDefaultBackoff(func() (*string, error) {
		if rand.Float32() > 0.9 {
			return nil, errors.New("simulated generator error")
		}
		if rand.Float32() > 0.5 {
			return nil, nil // simulated empty poll
		}
		if pos == len(returns) {
			return nil, reactive.GeneratorFinished{}
		}
		ret := returns[pos]
		pos++
		return &ret, nil
	})
	var response []string
	pipe.Observe(func(item string) error {
		response = append(response, item)
		return nil
	})
	pipe.Observe(func(item string) error {
		return errors.New("simulated sink error")
	})

	reactive.Async[string](pipe).Observe(func(s string) error {
		delay := time.Duration(rand.Intn(1000)) * time.Millisecond
		fmt.Println("Simulating long processing times", delay)

		time.Sleep(delay)
		return nil
	})
	pipe.Start()
	pipe.UponClose(func() {
		wg.Done()
	})
	wg.Wait()
	fmt.Println(response)
}
