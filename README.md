# Gotmr
[![Go Reference](https://pkg.go.dev/badge/github.com/livingsilver94/go-toml.svg)](https://pkg.go.dev/github.com/livingsilver94/gotmr) [![Go Report Card](https://goreportcard.com/badge/github.com/livingsilver94/gotmr)](https://goreportcard.com/report/github.com/livingsilver94/gotmr) ![GitHub](https://img.shields.io/github/license/livingsilver94/gotmr?color=red)

The Gotmr module provides a thread-safe and easy-to-use Go `time.Timer` wrapper.

`time.Timer` is notoriously hard to use in concurrent programs, because developers must ensure the internal channel is drained when Stopping and/or Resetting the timer, even in the case they don't know if there are goroutines listening to the internal channel. Gotmr's `Timer` hides this complexity to developers, which now may call Stop and Reset intuitively.

## Functioning
Gotmr's Timer is thread-safe because there's only one goroutine in control of it, while the ticks and the commands (Stop, Reset) are transferred using channels. Those channels are totally transparent to the user. Also, since channels are inherently thread-safe, there is no custom logic to handle concurrency inside the `Timer`'s code. The wrapper around `time.Timer` is very, very thin.
