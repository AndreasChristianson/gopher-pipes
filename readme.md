# Gopher Pipes
[![main-build](https://github.com/AndreasChristianson/gopher-pipes/actions/workflows/main-build.yaml/badge.svg)](https://github.com/AndreasChristianson/gopher-pipes/actions/workflows/main-build.yaml)
![coverage](https://raw.githubusercontent.com/AndreasChristianson/gopher-pipes/badges/.badges/main/coverage.svg)

Simple source/sink abstraction around generator functions, channels, and observing.

## Import

```shell
go get github.com/AndreasChristianson/gopher-pipes
```

## Usage

### simple

```shell
import (
	"github.com/AndreasChristianson/gopher-pipes/reactive"
)

  
underTest := Just("Hello")
```
