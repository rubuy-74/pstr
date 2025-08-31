package state_machine

import (
	"github.com/rubuy-74/pstr/internal/models/state"
	"github.com/rubuy-74/pstr/internal/parser"
	"github.com/rubuy-74/pstr/internal/utils"
)

func ToNFA(ctx *parser.ParseContext) *state.State {
	startOld, endOld := (ctx.Tokens[0]).ToNFA()
	for i := 1; i < len(ctx.Tokens); i++ {
		startNew, endNew := (ctx.Tokens[i]).ToNFA()
		endOld.Transitions[utils.Epsilon] = append(endOld.Transitions[utils.Epsilon], startNew)
		endOld = endNew
	}

	initialGlobalState := &state.State{
		Initial: true,
		Final:   false,
		Transitions: map[uint8][]*state.State{
			utils.Epsilon: {startOld},
		},
	}

	finalGlobalState := &state.State{
		Initial:     false,
		Final:       true,
		Transitions: map[uint8][]*state.State{},
	}

	endOld.Transitions[utils.Epsilon] = append(endOld.Transitions[utils.Epsilon], finalGlobalState)

	return initialGlobalState
}
