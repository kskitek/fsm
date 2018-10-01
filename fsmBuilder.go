package fsm

import (
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type EmitterDecorator func(Emitter) Emitter

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

func (f *FsmBuilder) WithErrorHandler(errorHandler ErrorHandler) *FsmBuilder {
	f.errorHandler = errorHandler
	return f
}

func (f *FsmBuilder) WithStateDecorator(s State, d EmitterDecorator) *FsmBuilder {
	e := f.stateDecorators[s]
	if e != nil {
		f.stateDecorators[s] = append(e, d)
	} else {
		f.stateDecorators[s] = []EmitterDecorator{d}
	}
	return f
}

func (f *FsmBuilder) WithDecorator(d EmitterDecorator) *FsmBuilder {
	f.decorators = append(f.decorators, d)
	return f
}

func defaultErrorHandler(curr State, err error) (State, error) {
	resErr := errors.Errorf("Error during transition from state {%s} of state: %s", string(curr), err.Error())
	log.Error(resErr)
	return willNotOccurState, resErr
}
