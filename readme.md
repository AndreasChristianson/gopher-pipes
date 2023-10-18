# Gopher Pipes

Simple source/sink abstraction around generator functions, channels, mapping, and observing.

### CI
[![coverage badge](https://github.com/AndreasChristianson/gopher-pipes/actions/workflows/coverage-badge.yaml/badge.svg)](https://github.com/AndreasChristianson/gopher-pipes/actions/workflows/coverage-badge.yaml)
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
  - [ ] coverage badge
- [ ] how should we handle errors?
  - [ ] in sinks
  - [x] in sources
- [ ] backpressure
- [x] retries
- [x] exponential backoff
  - [ ] configurable (partially complete)
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
- [ ] ci
  - [ ] code coverage
  - [ ] cut versions
- [ ] document
  - [ ] get us into golang docs
  - [ ] comments with links
