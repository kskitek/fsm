package fsm

import (
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

func Test_Error(t *testing.T) {
	v := &verifyableHandler{error: nil, state: Test2}
	m := map[State]Emitter{Test1: v.ToState2}
	out := New(m).Build()

	out.Start(Test1)
	actual := out.GetCurrent()

	assert.Equal(t, Test2, actual)
	assert.True(t, v.wasCalled)
}
