package fsm

import (
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

// EmitterDecorator is used by FsmBuilder to decorate both, individual step Emitter and global decorator.
type EmitterDecorator func(Emitter) Emitter

// New creates new FsmBuilder.
func New(handlers map[State]Emitter) *FsmBuilder {
	return &FsmBuilder{
		handlers:        handlers,
		errorHandler:    defaultErrorHandler,
		stateDecorators: make(map[State][]EmitterDecorator),
		decorators:      make([]EmitterDecorator, 0),
	}
}

type FsmBuilder struct {
	handlers        map[State]Emitter
	errorHandler    ErrorHandler
	stateDecorators map[State][]EmitterDecorator
	decorators      []EmitterDecorator
}

// Build prepares ready to start Fsm.
func (f *FsmBuilder) Build() Fsm {
	f.buildStateDecorators()
	f.buildDecorators()
	return &fsm{
		handlers:   f.handlers,
		errHandler: f.errorHandler,
	}
}

func (f *FsmBuilder) buildStateDecorators() {
	for state, decorators := range f.stateDecorators {
		for _, d := range decorators {
			fun := f.handlers[state]
			f.handlers[state] = d(fun)
		}
	}
}

func (f *FsmBuilder) buildDecorators() {
	for state, fun := range f.handlers {
		for _, d := range f.decorators {
			f.handlers[state] = d(fun)
		}
	}
}

// WithErrorHandler allows to handle errors in state Emitters. When none specified, defaultErrorHandler will be used.
func (f *FsmBuilder) WithErrorHandler(errorHandler ErrorHandler) *FsmBuilder {
	f.errorHandler = errorHandler
	return f
}

// WithStateDecorator adds EmitterDecorator only to given state.
// State Decorators will be applied in order in which they were declared.
//
// State Decorators are applied before global decorators.
func (f *FsmBuilder) WithStateDecorator(s State, d EmitterDecorator) *FsmBuilder {
	e := f.stateDecorators[s]
	if e != nil {
		f.stateDecorators[s] = append(e, d)
	} else {
		f.stateDecorators[s] = []EmitterDecorator{d}
	}
	return f
}

// WithDecorator adds EmitterDecorator global decorator to all states.
// Decorators will be applied in order in which they were declared.
//
// Decorators are applied after state decorators.
func (f *FsmBuilder) WithDecorator(d EmitterDecorator) *FsmBuilder {
	f.decorators = append(f.decorators, d)
	return f
}

func defaultErrorHandler(curr State, err error) (State, error) {
	resErr := errors.Errorf("Error during transition from state {%s} of state: %s", string(curr), err.Error())
	log.Error(resErr)
	return willNotOccurState, resErr
}
