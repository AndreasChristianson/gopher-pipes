# Gopher Pipes

Simple source/sink abstraction around generator functions, channels, mapping, and observing.

### CI
[![Go](https://github.com/AndreasChristianson/gopher-pipes/actions/workflows/test.yaml/badge.svg?branch=main&event=push)](https://github.com/AndreasChristianson/gopher-pipes/actions/workflows/test.yaml)


todo:

- [ ] readme
  - [ ] examples
  - [ ] getting started
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
