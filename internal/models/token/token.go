package token

import (
	"fmt"

	"github.com/rubuy-74/pstr/internal/models/state"
	"github.com/rubuy-74/pstr/internal/models/token_type"
	"github.com/rubuy-74/pstr/internal/utils"
)

type Token struct {
	TokenType token_type.TokenType
	Value     any
}

func (t Token) String() string {
	if b, ok := t.Value.(byte); ok {
		return fmt.Sprintf("{ %v '%v' }", t.TokenType, string(byte(b)))
	}
	return fmt.Sprintf("{ %v %v }", t.TokenType, t.Value)
}

type RepeatPayload struct {
	Min   int
	Max   int
	Token Token
}

func (rp RepeatPayload) String() string {
	return fmt.Sprintf("{ %v %v %v }", rp.Min, rp.Max, rp.Token.String())
}

type BracketPayload struct {
	Begin byte
	End   byte
}

func (bp BracketPayload) String() string {
	return fmt.Sprintf("[ '%v' - '%v' ]", string(bp.Begin), string(bp.End))
}

func (token Token) ToNFA() (*state.State, *state.State) {
	start := &state.State{
		Transitions: map[uint8][]*state.State{},
	}
	end := &state.State{
		Transitions: map[uint8][]*state.State{},
	}

	switch token.TokenType {
	case token_type.Group, token_type.GroupUncaptured: // UNTESTED
		values := token.Value.([]Token)
		start, end = (values[0]).ToNFA()
		for i := 1; i < len(values); i++ {
			startNew, endNew := (values[i]).ToNFA()
			end.Transitions[utils.Epsilon] = append(
				end.Transitions[utils.Epsilon],
				startNew,
			)
			end = endNew
		}
	case token_type.Bracket:
		values := token.Value.([]BracketPayload)
		for _, bp := range values {
			for i := bp.Begin; i < bp.End+1; i++ {
				fmt.Printf("%v ", i)
				start.Transitions[i] = []*state.State{end}
			}
		}
		fmt.Println()
	case token_type.Or:
		values := token.Value.([]Token)
		left := values[0]
		right := values[1]
		s1, e1 := left.ToNFA()
		s2, e2 := right.ToNFA()

		start.Transitions[utils.Epsilon] = []*state.State{s1, s2}
		e1.Transitions[utils.Epsilon] = []*state.State{end}
		e2.Transitions[utils.Epsilon] = []*state.State{end}

	case token_type.Repeat:
		payload := token.Value.(RepeatPayload)

		if payload.Min == 0 {
			start.Transitions[utils.Epsilon] = []*state.State{end}
		}

		copyCount := 0

		if payload.Max == utils.Infinite {
			if payload.Min == 0 {
				copyCount = 1
			} else {
				copyCount = payload.Min
			}
		} else {
			copyCount = payload.Max
		}

		sOld, eOld := payload.Token.ToNFA()
		start.Transitions[utils.Epsilon] = append(
			start.Transitions[utils.Epsilon],
			sOld,
		)

		for i := 2; i <= copyCount; i++ {
			sNew, eNew := payload.Token.ToNFA()

			eOld.Transitions[utils.Epsilon] = append(
				eOld.Transitions[utils.Epsilon],
				sNew,
			)

			sOld = sNew
			eOld = eNew

			// TODO: remove optional to improve performance
			if i > payload.Min {
				sNew.Transitions[utils.Epsilon] = append(
					sNew.Transitions[utils.Epsilon],
					end,
				)
			}
		}

		eOld.Transitions[utils.Epsilon] = append(
			eOld.Transitions[utils.Epsilon],
			end,
		)

		if payload.Max == utils.Infinite {
			end.Transitions[utils.Epsilon] = append(
				end.Transitions[utils.Epsilon],
				sOld,
			)
		}

	case token_type.Literal:
		ch := token.Value.(uint8)
		start.Transitions[ch] = []*state.State{end}

	default:
		panic("unknown type of token")
	}

	return start, end
}
