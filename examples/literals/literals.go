package main

import (
	"fmt"
	"github.com/AndreasChristianson/gopher-pipes/reactive"
)

func main() {
	pipe := reactive.Just("Hello", " ", "world", "!\n")
	pipe.Observe(func(item string) error {
		fmt.Print(item)
		return nil
	})
	pipe.Start()
	pipe.AwaitCompletion()
}
