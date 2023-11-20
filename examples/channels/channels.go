package main

import (
	"fmt"
	"github.com/AndreasChristianson/gopher-pipes/reactive"
)

func main() {
	channel := make(chan string)
	pipe := reactive.FromChan(channel)
	pipe.Observe(func(item string) error {
		fmt.Print(item)
		return nil
	})
	pipe.Start()
	channel <- "Hello"
	channel <- " "
	channel <- "world"
	channel <- "!\n"
	close(channel)
}
