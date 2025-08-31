package state

import (
	"github.com/rubuy-74/pstr/internal/utils"
)

type State struct {
	Initial     bool
	Final       bool
	Transitions map[uint8][]*State
}

// TODO: use multithreading for more performance
func (s *State) Check(input string, pos int) bool {
	ch := utils.GetChar(input, pos)

	if ch == utils.EndOfText && s.Final {
		return true
	}

	if states := s.Transitions[ch]; len(states) > 0 {
		nextState := states[0]
		if nextState.Check(input, pos+1) {
			return true
		}
	}

	for _, state := range s.Transitions[utils.Epsilon] {
		if state.Check(input, pos) {
			return true
		}

		if ch == utils.StartOfText && state.Check(input, pos+1) {
			return true
		}
	}

	return false
}
