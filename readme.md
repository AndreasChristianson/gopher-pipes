# Gopher Pipes

Simple source/sink abstraction around generator functions, channels, mapping, and observing.

### CI
[![main-build](https://github.com/AndreasChristianson/gopher-pipes/actions/workflows/main-build.yaml/badge.svg)](https://github.com/AndreasChristianson/gopher-pipes/actions/workflows/main-build.yaml)
![coverage](https://raw.githubusercontent.com/AndreasChristianson/gopher-pipes/badges/.badges/main/coverage.svg)

todo:

- [ ] features
  - [x] map
  - [x] buffer
  - [x] generator
  - [ ] reduce
  - [ ] peek/tap
- [ ] readme
  - [ ] examples
  - [ ] getting started
  - [x] coverage badge
- [ ] how should we handle errors?
  - [ ] in sinks
  - [x] in sources
- [ ] backpressure
- [x] retries
- [x] exponential backoff
  - [x] configurable
- [x] only register a sink once
- [x] hard errors vs soft errors
- [x] don't start right away
  - actually starting right away is better
- [ ] pass a context
- [ ] logging
  - [ ] option to squelch
  - [ ] option to collect till pipe completes
  - [ ] option to record metrics
  - [ ] option for verbose
- [x] ci
  - [x] code coverage
  - [x] cut versions
- [ ] document
  - [x] get us into golang docs
  - [x] comments with links
