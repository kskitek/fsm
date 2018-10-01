package fsm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	Test1 State = iota
	Test2
	Test3
)

type verifyableHandler struct {
	calledTimes int
	state       State
	error       error
}

func (v *verifyableHandler) ToState(s State) (State, error) {
	v.calledTimes++
	return v.state, v.error
}

func (v *verifyableHandler) Err(State, error) (State, error) {
	v.calledTimes++
	return v.state, v.error
}

func (v *verifyableHandler) Decorator(e Emitter) Emitter {
	v.calledTimes++
	return e
}

func Test_InitialStateIsSet(t *testing.T) {
	out := New(make(map[State]Emitter)).Build()

	out.Start(Test1)
	actual := out.GetCurrent()
	assert.Equal(t, Test1, actual)
}

func Test_NoHandlers_InitialStateIsReached(t *testing.T) {
	out := New(make(map[State]Emitter)).Build()

	out.Start(Test1)
	actual := out.GetCurrent()
	assert.Equal(t, Test1, actual)
}

func Test_TransitionFunctionIsCalled(t *testing.T) {
	v := &verifyableHandler{error: nil, state: Test2}
	m := map[State]Emitter{Test1: v.ToState}
	out := New(m).Build()

	out.Start(Test1)
	actual := out.GetCurrent()

	assert.Equal(t, Test2, actual)
	assert.Equal(t, 1, v.calledTimes)
}

func Test_Error_IsHandledByErrorHandler(t *testing.T) {
	err := fmt.Errorf("test error")
	v := &verifyableHandler{error: err, state: Test2}
	m := map[State]Emitter{Test1: v.ToState}
	errH := &verifyableHandler{error: err}
	out := New(m).WithErrorHandler(errH.Err).Build()

	out.Start(Test1)
	actual := out.GetCurrent()

	assert.Equal(t, Test1, actual)
	assert.Equal(t, 1, errH.calledTimes)
}

func Test_NoErrorHandler_Error_DefaultErrorHandlerIsProvided(t *testing.T) {
	v := &verifyableHandler{error: fmt.Errorf("test error")}
	m := map[State]Emitter{Test1: v.ToState}
	out := New(m).Build()

	out.Start(Test1)
	actual := out.GetCurrent()

	assert.Equal(t, Test1, actual)
}

func Test_Error_ErrorHandlerChangesState(t *testing.T) {
	v := &verifyableHandler{error: fmt.Errorf("test error"), state: Test2}
	m := map[State]Emitter{Test1: v.ToState}
	errH := &verifyableHandler{state: Test3}
	out := New(m).WithErrorHandler(errH.Err).Build()

	out.Start(Test1)
	actual := out.GetCurrent()

	assert.Equal(t, Test3, actual)
	assert.Equal(t, 1, errH.calledTimes)
}

func Test_Decorator_IsCalledForEachState(t *testing.T) {
	to2 := &verifyableHandler{state: Test2}
	to3 := &verifyableHandler{state: Test3}
	m := map[State]Emitter{
		Test1: to2.ToState,
		Test2: to3.ToState,
	}
	decorator := &verifyableHandler{state: Test3}

	out := New(m).WithDecorator(decorator.Decorator).Build()

	out.Start(Test1)
	actual := out.GetCurrent()

	assert.Equal(t, Test3, actual)
	assert.Equal(t, 2, decorator.calledTimes)
}

func Test_StateDecorator_IsCalledOnlyForState(t *testing.T) {
	to2 := &verifyableHandler{state: Test2}
	to3 := &verifyableHandler{state: Test3}
	m := map[State]Emitter{
		Test1: to2.ToState,
		Test2: to3.ToState,
	}
	decorator1 := &verifyableHandler{state: Test3}
	decorator2 := &verifyableHandler{state: Test3}

	out := New(m).
		WithStateDecorator(Test2, decorator1.Decorator).
		WithStateDecorator(Test2, decorator2.Decorator).
		Build()

	out.Start(Test1)
	actual := out.GetCurrent()

	assert.Equal(t, Test3, actual)
	assert.Equal(t, 1, decorator1.calledTimes)
	assert.Equal(t, 1, decorator2.calledTimes)

}
