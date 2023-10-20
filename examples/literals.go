package main

import (
	"github.com/AndreasChristianson/gopher-pipes/reactive"
	"time"
)

func main() {
	pipe := reactive.Just("Hello", " ", "world")
	pipe.UponClose(func() {
		println("!")
	})
	pipe.Observe(func(item string) error {
		print(item)
		return nil
	})
	<-time.After(time.Millisecond)
}
