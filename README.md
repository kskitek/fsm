[![GolangCI](https://golangci.com/badges/github.com/kskitek/fsm.svg)](https://golangci.com/r/github.com/kskitek/fsm) [![Go Report Card](https://goreportcard.com/badge/github.com/kskitek/fsm)](https://goreportcard.com/report/github.com/kskitek/fsm)

# fsm - Finite State Machine

`fsm` is very simple finite state machine implementation. The FSM is defined by States and next State Emitters.

## Usage

```go
package main

import "github.com/kskitek/fsm"

const (
    State1 fsm.State = iota
    State2
)

func main() {
    m := map[fsm.State]fsm.Emitter{State1: ToState2}
    
    out := fsm.New(m).Build()

    // start with initial state
    out.Start(State1)
}
    
func ToState2(s fsm.State) (fsm.State, error) {
    return State2, nil
}
```

## Decorators

`fsm` provides `fsm.EmitterDecorator` type that allows to wrap around Emitters, change state, log transitions,
send or store events, etc. You can apply decorator to either specified State or to all States using FsmBuilder.

**Example usage:**

```go
func Logging(e Emitter) Emitter {
	return func(s State) (State, error) {
		log.Printf("Started: %s\n", string(s))
		next, err := e(s)
		log.Printf("Ended: %s\n", string(s))
		return next, err
	}
}
```
