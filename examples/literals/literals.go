package main

import (
	"fmt"
	"github.com/AndreasChristianson/gopher-pipes/reactive"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	pipe := reactive.Just("Hello", " ", "world", "!\n")
	pipe.UponClose(func() {
		wg.Done()
	})
	pipe.Observe(func(item string) error {
		fmt.Print(item)
		return nil
	})
	pipe.Start()
	wg.Wait()
}
