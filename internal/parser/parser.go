package parser

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const infinite = -1

type TokenType uint8
type rangeSize int

const (
	group           TokenType = iota
	bracket         TokenType = iota
	or              TokenType = iota
	repeat          TokenType = iota
	literal         TokenType = iota
	groupUncaptured TokenType = iota
)

func (t TokenType) String() string {
	switch t {
	case group:
		return "group"
	case bracket:
		return "bracket"
	case or:
		return "or"
	case repeat:
		return "repeat"
	case literal:
		return "literal"
	case groupUncaptured:
		return "groupUncaptured"
	default:
		return fmt.Sprintf("TokenType(%d)", t)
	}
}

const (
	minInfinite rangeSize = iota
	InfiniteMax rangeSize = iota
	minMax      rangeSize = iota
)

type Token struct {
	TokenType TokenType
	Value     any
}

func (t Token) String() string {
	if b, ok := t.Value.(byte); ok {
		return fmt.Sprintf("{ %v '%v' }", t.TokenType, string(byte(b)))
	}
	return fmt.Sprintf("{ %v %v }", t.TokenType, t.Value)
}

type repeatPayload struct {
	min   int
	max   int
	token Token
}

func (rp repeatPayload) String() string {
	return fmt.Sprintf("{ %v %v %v }", rp.min, rp.max, rp.token.String())
}

type ParseContext struct {
	Pos    int
	Tokens []Token
}

type bracketPayload struct {
	begin byte
	end   byte
}

func (bp bracketPayload) String() string {
	return fmt.Sprintf("[ '%v' - '%v' ]", string(bp.begin), string(bp.end))
}

func (ctx ParseContext) Print() {
	fmt.Printf("ctx.Pos		: %v\n", ctx.Pos)
	fmt.Printf("ctx.Tokens: %v\n", ctx.Tokens)
}

/*
process parses one character at the current position of the regex string
and updates the ParseContext with tokens according to the symbol:
- '(' : start of a capturing group → delegates to processGroup
- '[' : start of a character class → delegates to processBrackets
- '|' : alternation (OR) → delegated to processOr (not yet implemented)
- '*' : repetition 0 or more times → processRepeat with min=0, max=infinite
- '+' : repetition 1 or more times → processRepeat with min=1, max=infinite
- '?' : repetition 0 or 1 → processRepeat with min=0, max=1
- '{' : repetition with explicit {min,max} → parses bounds with getMinMaxRange
- default: any other character is treated as a literal token
*/
func process(regex []byte, ctx *ParseContext) error {
	ch := regex[ctx.Pos]
	switch ch {
	case '(':
		err := processGroup(regex, ctx)
		if err != nil {
			return err
		}
	case '[':
		err := processBrackets(regex, ctx)
		if err != nil {
			return err
		}
	case '|': // TODO: not implemented
		processOr(regex, ctx)

	case '*':
		processRepeat(regex, ctx, 0, infinite)

	case '+':
		processRepeat(regex, ctx, 1, infinite)

	case '?':
		processRepeat(regex, ctx, 0, 1)

	case '{':
		minimum, maximum := getMinMaxRange(regex, ctx)
		processRepeat(regex, ctx, minimum, maximum)
	default:
		ctx.Tokens = append(ctx.Tokens,
			Token{
				TokenType: literal,
				Value:     ch,
			})
	}
	return nil
}

/*
getMinMaxRange extracts the min and max values from a repetition
range {m}, {m,}, or {m,n}.
- {m} → fixed repetition count
- {m,} → min repetitions with no upper bound
- {m,n} → explicit min and max repetitions
If parsing fails or invalid, returns infinite for both.
TODO: error handling
*/
func getMinMaxRange(regex []byte, ctx *ParseContext) (minimum int, maximum int) {
	newPos := findNextSymbol(regex, ctx.Pos, '}')
	rawRange := string(regex[ctx.Pos+1 : newPos])
	rangeString := strings.FieldsFunc(rawRange, func(r rune) bool {
		return r == ','
	})

	ctx.Pos = newPos
	// TODO: better checking (use macros)
	if !strings.Contains(rawRange, ",") {
		value, _ := strconv.Atoi(rawRange)
		return value, value
	} else {
		if len(rangeString) == 1 {
			if rawRange[0] == ',' {
				maximum, _ := strconv.Atoi(string(rangeString[0]))
				return infinite, maximum
			} else {
				minimum, _ := strconv.Atoi(string(rangeString[0]))
				return minimum, infinite
			}
		}
		if len(rangeString) == 2 {
			minimum, _ := strconv.Atoi(string(rangeString[0]))
			maximum, _ := strconv.Atoi(string(rangeString[1]))
			return minimum, maximum
		}
	}
	return infinite, infinite
}

/*
findNextSymbol scans forward in the regex string from prevPos
until it finds the specified symbol, and returns its index.
Returns -1 if not found within the regex.
*/
func findNextSymbol(regex []byte, prevPos int, symbol uint8) int {
	currPos := prevPos
	for regex[currPos] != symbol {
		currPos++
		if currPos > len(regex) {
			return -1
		}
	}
	return currPos
}

/*
chunkBytes splits a byte slice into chunks of given size,
but only keeps the first and last byte of each chunk.
Used mainly for bracket expressions with ranges (like a-z).
*/
func chunkBytes(data []byte, size int) []bracketPayload {
	numSlices := (len(data) + size - 1) / size
	subslices := make([]bracketPayload, 0, numSlices)

	for i := 0; i < len(data); i += size {
		end := min(i+size, len(data))
		subslice := bracketPayload{
			begin: data[i],
			end:   data[end-1],
		}
		subslices = append(subslices, subslice)
	}

	return subslices
}

/*
processGroup handles a capturing group "( ... )".
- Finds the closing ')'
- Recursively processes the inner substring as a new ParseContext
- Appends parsed tokens of the group back into the parent context
*/
func processGroup(regex []byte, ctx *ParseContext) error {
	ctx.Pos++
	newPos := findNextSymbol(regex, ctx.Pos, ')')
	if newPos == 1 {
		return fmt.Errorf("invalid ( in the regex string")
	}
	groupRegex := regex[ctx.Pos:newPos]
	groupCtx := &ParseContext{
		Pos:    0,
		Tokens: []Token{},
	}

	for groupCtx.Pos < len(groupRegex) {
		err := process(groupRegex, groupCtx)
		if err != nil {
			return err
		}
		groupCtx.Pos++
	}

	ctx.Pos = newPos
	ctx.Tokens = append(ctx.Tokens, groupCtx.Tokens...)

	return nil
}

/*
processBrackets handles a character class "[ ... ]".
- Finds the closing ']'
- If it contains '-', splits into ranges (e.g. a-z)
- Otherwise, treats as a single range between first and last characters
- Appends bracket tokens to the context
*/
func processBrackets(regex []byte, ctx *ParseContext) error {
	ctx.Pos++
	newPos := findNextSymbol(regex, ctx.Pos, ']')
	if newPos == 1 {
		return fmt.Errorf("invalid [ in the regex string")
	}
	insideRegex := regex[ctx.Pos:newPos]

	bpSlice := []bracketPayload{}

	if slices.Contains(insideRegex, '-') {
		ranges := chunkBytes(insideRegex, 3)
		for _, r := range ranges {
			bpSlice = append(bpSlice, r)
		}
	} else {
		bpSlice = append(
			bpSlice,
			bracketPayload{
				begin: insideRegex[1],
				end:   insideRegex[len(insideRegex)-2],
			},
		)
	}

	token := Token{
		TokenType: bracket,
		Value:     bpSlice,
	}
	ctx.Tokens = append(ctx.Tokens, token)

	ctx.Pos = newPos
	return nil
}

/*
processOr placeholder for handling alternation '|'.
Currently not implemented.
*/
func processOr(regex []byte, ctx *ParseContext) {
	rhsContext := &ParseContext{
		Pos:    ctx.Pos,
		Tokens: []Token{},
	}

	rhsContext.Pos += 1
	for rhsContext.Pos < len(regex) && regex[rhsContext.Pos] != ')' {
		process(regex, rhsContext)
		rhsContext.Pos += 1
	}

	left := Token{
		TokenType: groupUncaptured,
		Value:     ctx.Tokens,
	}

	right := Token{
		TokenType: groupUncaptured,
		Value:     rhsContext.Tokens,
	}

	ctx.Pos = rhsContext.Pos

	ctx.Tokens = []Token{{
		TokenType: or,
		Value:     []Token{left, right},
	}}
}

/*
processRepeat placeholder for handling repetition operators (*, +, ?, {m,n}).
Currently not implemented.
- The current processRepeat only allows one token to be repeated
- It would be nice to be able to repeat a group. ex.: ([a-z]){2}
*/
func processRepeat(regex []byte, ctx *ParseContext, min int, max int) {
	_ = regex // TODO: the regex variable will be used in the future
	lastToken := ctx.Tokens[len(ctx.Tokens)-1]
	ctx.Tokens = ctx.Tokens[:len(ctx.Tokens)-1]
	ctx.Tokens = append(ctx.Tokens, Token{
		TokenType: repeat,
		Value: repeatPayload{
			min:   min,
			max:   max,
			token: lastToken,
		},
	})
}

/*
Parse initializes parsing for a regex string.
- Converts string to byte slice
- Iterates through each character, delegating to process
- Returns final ParseContext with tokens
*/
func Parse(regexString string) (*ParseContext, error) {
	regex := []byte(regexString)
	ctx := &ParseContext{
		Pos:    0,
		Tokens: []Token{},
	}
	for ctx.Pos < len(regex) {
		err := process(regex, ctx)
		if err != nil {
			return nil, err
		}
		ctx.Pos++
	}

	return ctx, nil
}
