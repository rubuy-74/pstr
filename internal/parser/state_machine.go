package parser

import "fmt"

const epsilon uint8 = 0

type state struct {
	Initial     bool
	Final       bool
	Transitions map[uint8][]*state
}

func tokenToNFA(token Token) (*state, *state) {
	start := &state{
		Transitions: map[uint8][]*state{},
	}
	end := &state{
		Transitions: map[uint8][]*state{},
	}

	switch token.TokenType {
	case group:
		values := token.Value.([]Token)
		start, end = tokenToNFA(values[0])
		for i := 1; i < len(values); i++ {
			startNew, endNew := tokenToNFA(values[i])
			end.Transitions[epsilon] = append(
				end.Transitions[epsilon],
				startNew,
			)
			end = endNew
		}
	case bracket:
		values := token.Value.([]bracketPayload)
		for _, bp := range values {
			for i := bp.begin; i < bp.end+1; i++ {
				fmt.Printf("%v ", i)
				start.Transitions[i] = []*state{end}
			}
		}
		fmt.Println()
	case or:
		values := token.Value.([]Token)
		left := values[0]
		right := values[1]
		s1, e1 := tokenToNFA(left)
		s2, e2 := tokenToNFA(right)

		start.Transitions[epsilon] = []*state{s1, s2}
		e1.Transitions[epsilon] = []*state{end}
		e2.Transitions[epsilon] = []*state{end}

	case repeat:
		payload := token.Value.(repeatPayload)

		if payload.min == 0 {
			start.Transitions[epsilon] = []*state{end}
		}

		copyCount := 0

		if payload.max == infinite {
			if payload.min == 0 {
				copyCount = 1
			} else {
				copyCount = payload.min
			}
		} else {
			copyCount = payload.max
		}

		sOld, eOld := tokenToNFA(payload.token)
		start.Transitions[epsilon] = append(
			start.Transitions[epsilon],
			sOld,
		)

		for i := 2; i <= copyCount; i++ {
			sNew, eNew := tokenToNFA(payload.token)

			eOld.Transitions[epsilon] = append(
				eOld.Transitions[epsilon],
				sNew,
			)

			sOld = sNew
			eOld = eNew

			// TODO: remove optional to improve performance
			if i > payload.min {
				sNew.Transitions[epsilon] = append(
					sNew.Transitions[epsilon],
					end,
				)
			}
		}

		eOld.Transitions[epsilon] = append(
			eOld.Transitions[epsilon],
			end,
		)

		if payload.max == infinite {
			end.Transitions[epsilon] = append(
				end.Transitions[epsilon],
				sOld,
			)
		}

	case literal:
		ch := token.Value.(uint8)
		start.Transitions[ch] = []*state{end}
	case groupUncaptured:
	default:
		panic("unknown type of token")
	}

	return start, end
}

func ToNFA(ctx *ParseContext) *state {
	startOld, endOld := tokenToNFA(ctx.Tokens[0])
	for i := 1; i < len(ctx.Tokens); i++ {
		startNew, endNew := tokenToNFA(ctx.Tokens[i])
		endOld.Transitions[epsilon] = append(endOld.Transitions[epsilon], startNew)
		endOld = endNew
	}

	initialGlobalState := &state{
		Initial: true,
		Final:   false,
		Transitions: map[uint8][]*state{
			epsilon: {startOld},
		},
	}

	finalGlobalState := &state{
		Initial:     false,
		Final:       true,
		Transitions: map[uint8][]*state{},
	}

	endOld.Transitions[epsilon] = append(endOld.Transitions[epsilon], finalGlobalState)

	return initialGlobalState
}
