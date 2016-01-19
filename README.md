# GPIO test tool

A simple tool that recreates functionality of
[elr-hw](https://github.com/open-rnd/erl-hw) using Go. The tool also
gives some insight into the performance of Go's scheduler and garbage
collector. In our case, we've written this code to compare how Go's
runtime fares when compared to Erlang. As such it may be treated as a
benchmark of sorts, however make sure to take the results with a grain
of salt.

## Contept

The idea is to measure how the system behaves when a large number of
processes is spawned. The do this the benchmark will create a couple
of chains, each consisting of processes. In each chain, the processes
will pass around a token from one to another and once a complete
revolution/cycle is complete, the state of assigned LED/GPIO will be
triggered.

In order to map the concepts from Erlang into Go, a goroutine is used
as a process and the token is passed over channels. The chain is
implemented inside `src/chain/chain.go`. To limit the effect of
spawning a large number of goroutines on the performance, we first
call `chain.Spawn()`, wait a little bit, and then call `chain.Start()`
to inject the token, thus starting the sequence.

Measurements are done on a tiny Linux system such as Raspberry Pi. A
breadboard with LEDs is connected to GPIO pins of the board. The
benchmark is expected to toggle the state of LEDs when running. For
accurate timing measurements, connect a logical state analyzer to the
GPIO to record a trace.

## Building

To build `blinktest` for Raspberry Pi run the following command:

```
GOARCH=arm GOARM=6 GOPATH=$PWD go build -v blinktest
```

Make sure to set `GOARM` appropriately when building for other ARM
based boards.

## Running

Before running on an embedded system make sure to adjust the GPIO/LED
mappings in `src/blinktest/blinktest.go`. The GPIOs are accessed via
`/sys/class/gpio` interface. In practice LEDs are not necessary,
however they provide a nice visual feedback.

## Raspberry Pi

The benchmark uses 4 GPIO pins: 4, 17, 27, 22.

## BeagleBone Black

Make sure to set `GOARM=7` when building. GPIO pins on P8 header: 67
(P8.08), 68 (P8.10), 44 (P8.12), 24 (P8.14).

The GPIO pins can be assigned multiple functions, you need to
configure them into GPIO mode before running the benchmark. This can
be done by running a `config-pin` helper tool like this:

```
for p in P8.08 P8.10 P8.12 P8.14; do config-pin $p gpio; done
```
