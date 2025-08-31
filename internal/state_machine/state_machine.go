package statemachine

import (
	"github.com/rubuy-74/pstr/internal/parser"
)

type tokenValue uint8

const epsilon tokenValue = 0

type state struct {
	Initial     bool
	Final       bool
	Transitions map[tokenValue][]*state
}

func tokenToNFA(token parser.Token) (*state, *state) {
	start := &state{
		Transitions: map[tokenValue][]*state{},
	}
	end := &state{
		Transitions: map[tokenValue][]*state{},
	}

	switch token.TokenType {
	case parser.group:
	case bracket:
	case or:
	case repeat:
	case literal:
	case groupUncaptured:

	}

	return start, end
}

func toNFA(ctx *parser.ParseContext) *state {
	startOld, endOld := tokenToNFA(ctx.Tokens[0])
	for i := 1; i < len(ctx.Tokens); i++ {
		startNew, endNew := tokenToNFA(ctx.Tokens[i])
		endOld.Transitions[epsilon] = append(endOld.Transitions[epsilon], startNew)
		endOld = endNew
	}

	initialGlobalState := &state{
		Initial: true,
		Final:   false,
		Transitions: map[tokenValue][]*state{
			epsilon: {startOld},
		},
	}

	finalGlobalState := &state{
		Initial:     false,
		Final:       true,
		Transitions: map[tokenValue][]*state{},
	}

	endOld.Transitions[epsilon] = append(endOld.Transitions[epsilon], finalGlobalState)

	return initialGlobalState
}
