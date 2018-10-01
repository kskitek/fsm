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
	wasCalled bool
	state     State
	error     error
}

func (v *verifyableHandler) ToState2() (State, error) {
	v.wasCalled = true
	return v.state, v.error
}

func (v *verifyableHandler) Err(State, error) (State, error) {
	v.wasCalled = true
	return v.state, v.error
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
	m := map[State]Emitter{Test1: v.ToState2}
	out := New(m).Build()

	out.Start(Test1)
	actual := out.GetCurrent()

	assert.Equal(t, Test2, actual)
	assert.True(t, v.wasCalled)
}

func Test_Error_IsHandledByErrorHandler(t *testing.T) {
	err := fmt.Errorf("test error")
	v := &verifyableHandler{error: err, state: Test2}
	m := map[State]Emitter{Test1: v.ToState2}
	errH := &verifyableHandler{error: err}
	out := New(m).WithErrorHandler(errH.Err).Build()

	out.Start(Test1)
	actual := out.GetCurrent()

	assert.Equal(t, Test1, actual)
	assert.True(t, errH.wasCalled)
}

func Test_Error_ErrorHandlerChangesState(t *testing.T) {
	v := &verifyableHandler{error: fmt.Errorf("test error"), state: Test2}
	m := map[State]Emitter{Test1: v.ToState2}
	errH := &verifyableHandler{state: Test3}
	out := New(m).WithErrorHandler(errH.Err).Build()

	out.Start(Test1)
	actual := out.GetCurrent()

	assert.Equal(t, Test3, actual)
}
