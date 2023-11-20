package main

import (
	"fmt"
	"github.com/AndreasChristianson/gopher-pipes/reactive"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	returns := []string{"Hello", " ", "world", "!\n"}
	pos := 0
	wg.Add(1)
	pipe := reactive.FromGenerator(func() (*string, error) {
		if pos == len(returns) {
			return nil, reactive.GeneratorFinished{}
		}
		ret := returns[pos]
		pos++
		return &ret, nil
	})
	pipe.Observe(func(item string) error {
		fmt.Print(item)
		return nil
	})
	pipe.UponClose(func() {
		wg.Done()
	})
	pipe.Start()
	wg.Wait()
}
