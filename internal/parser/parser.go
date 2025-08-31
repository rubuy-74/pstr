package parser

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const infinite = -1

type tokenType uint8
type rangeSize int

const (
	group           tokenType = iota
	bracket         tokenType = iota
	or              tokenType = iota
	repeat          tokenType = iota
	literal         tokenType = iota
	groupUncaptured tokenType = iota
)

func (t tokenType) String() string {
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
		return fmt.Sprintf("tokenType(%d)", t)
	}
}

const (
	minInfinite rangeSize = iota
	InfiniteMax rangeSize = iota
	minMax      rangeSize = iota
)

type token struct {
	tokenType tokenType
	value     any
}

func (t token) String() string {
	if b, ok := t.value.(byte); ok {
		return fmt.Sprintf("{ %v '%v' }", t.tokenType, string(byte(b)))
	}
	return fmt.Sprintf("{ %v %v }", t.tokenType, t.value)
}

type repeatPayload struct {
	min   int
	max   int
	token token
}

func (rp repeatPayload) String() string {
	return fmt.Sprintf("{ %v %v %v }", rp.min, rp.max, rp.token.String())
}

type ParseContext struct {
	pos    int
	tokens []token
}

type bracketPayload struct {
	begin byte
	end   byte
}

func (bp bracketPayload) String() string {
	return fmt.Sprintf("[ '%v' - '%v' ]", string(bp.begin), string(bp.end))
}

func (ctx ParseContext) Print() {
	fmt.Printf("ctx.pos		: %v\n", ctx.pos)
	fmt.Printf("ctx.tokens: %v\n", ctx.tokens)
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
	ch := regex[ctx.pos]
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
		ctx.tokens = append(ctx.tokens,
			token{
				tokenType: literal,
				value:     ch,
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
	newPos := findNextSymbol(regex, ctx.pos, '}')
	rawRange := string(regex[ctx.pos+1 : newPos])
	rangeString := strings.FieldsFunc(rawRange, func(r rune) bool {
		return r == ','
	})

	ctx.pos = newPos
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
	ctx.pos++
	newPos := findNextSymbol(regex, ctx.pos, ')')
	if newPos == 1 {
		return fmt.Errorf("invalid ( in the regex string")
	}
	groupRegex := regex[ctx.pos:newPos]
	groupCtx := &ParseContext{
		pos:    0,
		tokens: []token{},
	}

	for groupCtx.pos < len(groupRegex) {
		err := process(groupRegex, groupCtx)
		if err != nil {
			return err
		}
		groupCtx.pos++
	}

	ctx.pos = newPos
	ctx.tokens = append(ctx.tokens, groupCtx.tokens...)

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
	ctx.pos++
	newPos := findNextSymbol(regex, ctx.pos, ']')
	if newPos == 1 {
		return fmt.Errorf("invalid [ in the regex string")
	}
	insideRegex := regex[ctx.pos:newPos]
	if slices.Contains(insideRegex, '-') {
		ranges := chunkBytes(insideRegex, 3)
		for _, r := range ranges {
			ctx.tokens = append(ctx.tokens, token{
				tokenType: bracket,
				value:     r,
			})
		}
	} else {
		ctx.tokens = append(ctx.tokens, token{
			tokenType: bracket,
			value: bracketPayload{
				begin: insideRegex[1],
				end:   insideRegex[len(insideRegex)-2],
			},
		})
	}

	ctx.pos = newPos
	return nil
}

/*
processOr placeholder for handling alternation '|'.
Currently not implemented.
*/
func processOr(regex []byte, ctx *ParseContext) {
	rhsContext := &ParseContext{
		pos:    ctx.pos,
		tokens: []token{},
	}

	rhsContext.pos += 1
	for rhsContext.pos < len(regex) && regex[rhsContext.pos] != ')' {
		process(regex, rhsContext)
		rhsContext.pos += 1
	}

	left := token{
		tokenType: groupUncaptured,
		value:     ctx.tokens,
	}

	right := token{
		tokenType: groupUncaptured,
		value:     rhsContext.tokens,
	}

	ctx.pos = rhsContext.pos

	ctx.tokens = []token{{
		tokenType: or,
		value:     []token{left, right},
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
	lastToken := ctx.tokens[len(ctx.tokens)-1]
	ctx.tokens = ctx.tokens[:len(ctx.tokens)-1]
	ctx.tokens = append(ctx.tokens, token{
		tokenType: repeat,
		value: repeatPayload{
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
		pos:    0,
		tokens: []token{},
	}
	for ctx.pos < len(regex) {
		err := process(regex, ctx)
		if err != nil {
			return nil, err
		}
		ctx.pos++
	}

	return ctx, nil
}
