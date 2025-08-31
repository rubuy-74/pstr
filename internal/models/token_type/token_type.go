package token_type

import "fmt"

type TokenType uint8

const (
	Group           TokenType = iota
	Bracket         TokenType = iota
	Or              TokenType = iota
	Repeat          TokenType = iota
	Literal         TokenType = iota
	GroupUncaptured TokenType = iota
)

func (t TokenType) String() string {
	switch t {
	case Group:
		return "group"
	case Bracket:
		return "bracket"
	case Or:
		return "or"
	case Repeat:
		return "repeat"
	case Literal:
		return "literal"
	case GroupUncaptured:
		return "groupUncaptured"
	default:
		return fmt.Sprintf("TokenType(%d)", t)
	}
}
