package parser

const (
	startOfText uint8 = 2
	endOfText   uint8 = 3
)

func getChar(input string, pos int) uint8 {
	if pos >= len(input) {
		return endOfText
	}

	if pos < 0 {
		return startOfText
	}

	return uint8(input[pos])
}

// TODO: use multithreading for more performance
func (s *state) Check(input string, pos int) bool {
	ch := getChar(input, pos)

	if ch == endOfText && s.Final {
		return true
	}

	if states := s.Transitions[ch]; len(states) > 0 {
		nextState := states[0]
		if nextState.Check(input, pos+1) {
			return true
		}
	}

	for _, state := range s.Transitions[epsilon] {
		if state.Check(input, pos) {
			return true
		}

		if ch == startOfText && state.Check(input, pos+1) {
			return true
		}
	}

	return false
}
