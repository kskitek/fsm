package fsm

import (
	"github.com/pkg/errors"
	"log"
)

func Logging(e Emitter) Emitter {
	return func(s State) (State, error) {
		log.Printf("Started: %s\n", string(s))
		next, err := e(s)
		log.Printf("Ended: %s\n", string(s))
		return next, err
	}
}

type RepeatingErrorHandler struct {
	count     int
	MaxCount  int
	lastState State
}

func (r *RepeatingErrorHandler) Err(s State, err error) (State, error) {
	if s == r.lastState {
		r.count++
		if r.count == r.MaxCount {
			return s, errors.Errorf("Maximum retries count reached when handling state {%s}", string(s))
		}
		return s, nil
	}
	r.lastState = s
	r.count = 0
	return s, nil
}
