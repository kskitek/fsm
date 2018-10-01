package fsm

type State int

type Emitter = func() (State, error)

type ErrorHandler = func(error) *State

// TODO add options with fluent API like WithErrorHandler.WithTransitionDecorator
// TODO New with stateDecorators
// TODO add errorHandler and default errorHandler
// TODO options like State and transition timeout, whole FSM timeout or state repeat count limit
// some of those can be implemented as transition stateDecorators

type Fsm interface {
	Start(State)
	GetCurrent() State
}

type fsm struct {
	state    State
	handlers map[State]Emitter
}

func (f *fsm) Start(initial State) {
	f.state = initial

	for {
		transition, ok := f.handlers[f.state]
		if !ok {
			return
		}
		next, err := transition()
		if err != nil {
			// TODO error handler
			return
		}
		f.state = next
	}
}

func (f *fsm) GetCurrent() State {
	return f.state
}
