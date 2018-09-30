package fsm

type State int

type Transition = func() (State, error)
type StateHandler struct {
	State          State
	TransitionFunc Transition
}

type ErrorHandler = func(error) *State

func SimpleTransition(from State, to State) *StateHandler {
	f := func() (State, error) { return to, nil }
	return &StateHandler{State: from, TransitionFunc: f}
}

// TODO constructor from map
func New(handlers []StateHandler) Fsm {
	hmap := make(map[State]Transition, len(handlers))
	for _, v := range handlers {
		hmap[v.State] = v.TransitionFunc
	}
	return NewFromMap(hmap)
}

func NewFromMap(handlers map[State]Transition) Fsm {
	return &fsm{
		handlers: handlers,
	}
}

// TODO add options with fluent API like WithErrorHandler.WithTransitionDecorator
// TODO New with decorators
// TODO add errorHandler and default errorHandler
// TODO options like State and transition timeout, whole FSM timeout or state repeat count limit
// some of those can be implemented as transition decorators

type Fsm interface {
	Start(State)
	GetCurrent() State
}

type fsm struct {
	state    State
	handlers map[State]Transition
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
