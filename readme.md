# Gopher Pipes
[![main-build](https://github.com/AndreasChristianson/gopher-pipes/actions/workflows/main-build.yaml/badge.svg)](https://github.com/AndreasChristianson/gopher-pipes/actions/workflows/main-build.yaml)
![coverage](https://raw.githubusercontent.com/AndreasChristianson/gopher-pipes/badges/.badges/main/coverage.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/AndreasChristianson/gopher-pipes.svg)](https://pkg.go.dev/github.com/AndreasChristianson/gopher-pipes/reactive)

Simple source/sink abstraction around generator functions, channels, and observing.

## Import

```shell
go get github.com/AndreasChristianson/gopher-pipes
```

## Usage

see [the examples folder](/examples) for more examples.

### simple

This example creates a source with four strings, observes them, and prints the observed strings.
Verbose logging is enabled.

```go
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
```

### complicated

This example polls a persistent redis stream using [XREAD](https://redis.io/commands/xread/)
(via [go-redis](https://github.com/redis/go-redis)) and routes messages to a [websocket](https://github.com/gorilla/websocket)

```go
// todo
```
