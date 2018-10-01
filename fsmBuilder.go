package fsm

import "github.com/prometheus/common/log"

func New(handlers map[State]Emitter) *FsmBuilder {
	return &FsmBuilder{
		handlers:        handlers,
		errorHandler:    defaultErrorHandler,
		stateDecorators: make(map[State][]Emitter),
		decorators:      make([]Emitter, 0),
	}
}

type FsmBuilder struct {
	handlers        map[State]Emitter
	errorHandler    ErrorHandler
	stateDecorators map[State][]Emitter
	decorators      []Emitter
}

func (f *FsmBuilder) Build() Fsm {
	return &fsm{
		handlers:   f.handlers,
		errHandler: f.errorHandler,
	}
}

func (f *FsmBuilder) WithErrorHandler(errorHandler ErrorHandler) *FsmBuilder {
	f.errorHandler = errorHandler
	return f
}

func (f *FsmBuilder) WithStateDecorator(s State, d Emitter) *FsmBuilder {
	e := f.stateDecorators[s]
	if e != nil {
		f.stateDecorators[s] = append(e, d)
	} else {
		f.stateDecorators[s] = []Emitter{d}
	}
	return f
}

func (f *FsmBuilder) WithDecorator(d Emitter) *FsmBuilder {
	f.decorators = append(f.decorators, d)
	return f
}

func defaultErrorHandler(curr State, err error) (State, error) {
	log.Errorf("Error during transition from state {%s} of state: %s", string(curr), err.Error())
	return curr, nil
}
