package fsm

import "math"

const (
	willNotOccurState State = math.MinInt32
)

// State of Fsm
type State int

// Emitter emits next state from current state.
// When it returns error, fsm will either recover with error handler or stop.
type Emitter = func(State) (State, error)

// ErrorHandler func is used by Fsm to handle Emitter's errors.
//
// ErrorHandler can alter next State by returning it.
//
// When error is returned, state is ignored and Fsm will stop.
type ErrorHandler = func(curr State, err error) (State, error)

type Fsm interface {
	// Start fsm using initial state.
	//
	// Fsm will continue as long as Emitter functions are found for current state.
	// When there is no Emitter registered to handle state or unrecoverable error was returned, Fsm stops.
	Start(initial State)
	// GetCurrent returns last good state. In case of error it returns last state that Fsm was before error.
	GetCurrent() State
}

type fsm struct {
	state      State
	handlers   map[State]Emitter
	errHandler ErrorHandler
}

func (f *fsm) Start(initial State) {
	f.state = initial

	for {
		transition, ok := f.handlers[f.state]
		if !ok {
			return
		}
		next, err := transition(f.state)
		if err != nil {
			newNext, err := f.errHandler(f.state, err)
			if err != nil {
				return
			}
			next = newNext
		}
		f.state = next
	}
}

func (f *fsm) GetCurrent() State {
	return f.state
}
