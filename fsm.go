package fsm

import "math"

const (
	willNotOccurState State = math.MinInt32
)

type State int

type Emitter = func(State) (State, error)

type ErrorHandler = func(curr State, err error) (State, error)

type Fsm interface {
	Start(State)
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
