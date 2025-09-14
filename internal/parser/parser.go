package parser

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/rubuy-74/pstr/internal/models/token"
	"github.com/rubuy-74/pstr/internal/models/token_type"
	"github.com/rubuy-74/pstr/internal/utils"
)

type rangeSize int

type ParseContext struct {
	Pos    int
	Tokens []token.Token
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
	case '|':
		err := processOr(regex, ctx)
		if err != nil {
			return err
		}

	case '*':
		err := processRepeat(regex, ctx, 0, utils.Infinite)
		if err != nil {
			return err
		}

	case '+':
		err := processRepeat(regex, ctx, 1, utils.Infinite)
		if err != nil {
			return err
		}

	case '?':
		err := processRepeat(regex, ctx, 0, 1)
		if err != nil {
			return err
		}

	case '{':
		minimum, maximum, err := getMinMaxRange(regex, ctx)
		if err != nil {
			return err
		}
		err = processRepeat(regex, ctx, minimum, maximum)
		if err != nil {
			return err
		}
	default:
		ctx.Tokens = append(ctx.Tokens,
			token.Token{
				TokenType: token_type.Literal,
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
*/
func getMinMaxRange(regex []byte, ctx *ParseContext) (minimum int, maximum int, err error) {
	newPos, err := findNextSymbol(regex, ctx.Pos, '}')
	if err != nil {
		return utils.Infinite, utils.Infinite, fmt.Errorf("missing ending '}'")
	}
	rawRange := string(regex[ctx.Pos+1 : newPos])
	rangeString := strings.FieldsFunc(rawRange, func(r rune) bool {
		return r == ','
	})

	ctx.Pos = newPos
	// TODO: better checking (use macros)
	if !strings.Contains(rawRange, ",") {
		value, _ := strconv.Atoi(rawRange)
		return value, value, nil
	} else {
		if len(rangeString) == 1 {
			if rawRange[0] == ',' {
				maximum, _ := strconv.Atoi(string(rangeString[0]))
				return utils.Infinite, maximum, nil
			} else {
				minimum, _ := strconv.Atoi(string(rangeString[0]))
				return minimum, utils.Infinite, nil
			}
		}
		if len(rangeString) == 2 {
			minimum, _ := strconv.Atoi(string(rangeString[0]))
			maximum, _ := strconv.Atoi(string(rangeString[1]))
			return minimum, maximum, nil
		}
	}
	return utils.Infinite, utils.Infinite, fmt.Errorf("invalid range syntax")
}

/*
findNextSymbol scans forward in the regex string from prevPos
until it finds the specified symbol, and returns its index.
Returns -1 if not found within the regex.
*/
func findNextSymbol(regex []byte, prevPos int, symbol uint8) (int, error) {
	currPos := prevPos
	if currPos >= len(regex) {
		return -1, fmt.Errorf("symbol not found")
	}

	for regex[currPos] != symbol {
		currPos++
		if currPos >= len(regex) {
			return -1, fmt.Errorf("symbol not found")
		}
	}
	return currPos, nil
}

/*
chunkBytes splits a byte slice into chunks of given size,
but only keeps the first and last byte of each chunk.
Used mainly for bracket expressions with ranges (like a-z).
*/
func chunkBytes(data []byte, size int) []token.BracketPayload {
	numSlices := (len(data) + size - 1) / size
	subslices := make([]token.BracketPayload, 0, numSlices)

	for i := 0; i < len(data); i += size {
		end := min(i+size, len(data))
		subslice := token.BracketPayload{
			Begin: data[i],
			End:   data[end-1],
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
	newPos, err := findNextSymbol(regex, ctx.Pos, ')')
	if err != nil {
		return err
	}
	if newPos == 1 {
		return fmt.Errorf("invalid ( in the regex string")
	}
	groupRegex := regex[ctx.Pos:newPos]
	if len(groupRegex) == 0 {
		return fmt.Errorf("empty group")
	}
	groupCtx := &ParseContext{
		Pos:    0,
		Tokens: []token.Token{},
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
	newPos, err := findNextSymbol(regex, ctx.Pos, ']')
	if err != nil {
		return err
	}
	insideRegex := regex[ctx.Pos:newPos]
	if len(insideRegex) == 0 {
		return fmt.Errorf("empty bracket expression")
	}

	if newPos == 1 || len(insideRegex) < 2 {
		return fmt.Errorf("invalid [ in the regex string")
	}

	bpSlice := []token.BracketPayload{}

	if slices.Contains(insideRegex, '-') {
		ranges := chunkBytes(insideRegex, 3)
		bpSlice = append(bpSlice, ranges...)
	} else {
		bpSlice = append(
			bpSlice,
			token.BracketPayload{
				Begin: insideRegex[1],
				End:   insideRegex[len(insideRegex)-2],
			},
		)
	}

	token := token.Token{
		TokenType: token_type.Bracket,
		Value:     bpSlice,
	}
	ctx.Tokens = append(ctx.Tokens, token)

	ctx.Pos = newPos
	return nil
}

/*
 */
func processOr(regex []byte, ctx *ParseContext) error {
	rhsContext := &ParseContext{
		Pos:    ctx.Pos,
		Tokens: []token.Token{},
	}
	if len(ctx.Tokens) == 0 {
		return fmt.Errorf("missing left operand")
	}

	if rhsContext.Pos+1 >= len(regex) {
		return fmt.Errorf("missing right operand")
	}

	rhsContext.Pos += 1
	for rhsContext.Pos < len(regex) && regex[rhsContext.Pos] != ')' {
		err := process(regex, rhsContext)
		if err != nil {
			return err
		}
		rhsContext.Pos += 1
	}

	if len(rhsContext.Tokens) == 0 {
		return fmt.Errorf("missing right operand")
	}

	left := token.Token{
		TokenType: token_type.GroupUncaptured,
		Value:     ctx.Tokens,
	}

	right := token.Token{
		TokenType: token_type.GroupUncaptured,
		Value:     rhsContext.Tokens,
	}

	ctx.Pos = rhsContext.Pos

	ctx.Tokens = []token.Token{{
		TokenType: token_type.Or,
		Value:     []token.Token{left, right},
	}}

	return nil
}

/*
processRepeat placeholder for handling repetition operators (*, +, ?, {m,n}).
Currently not implemented.
- The current processRepeat only allows one token to be repeated
- It would be nice to be able to repeat a group. ex.: ([a-z]){2}
*/
func processRepeat(regex []byte, ctx *ParseContext, min int, max int) error {
	_ = regex // TODO: the regex variable will be used in the future

	if len(ctx.Tokens) == 0 {
		return fmt.Errorf("missing repeating element")
	}

	lastToken := ctx.Tokens[len(ctx.Tokens)-1]
	ctx.Tokens = ctx.Tokens[:len(ctx.Tokens)-1]
	ctx.Tokens = append(ctx.Tokens, token.Token{
		TokenType: token_type.Repeat,
		Value: token.RepeatPayload{
			Min:   min,
			Max:   max,
			Token: lastToken,
		},
	})

	return nil
}

/*
Parse initializes parsing for a regex string.
- Converts string to byte slice
- Iterates through each character, delegating to process
- Returns final ParseContext with tokens
*/
func Parse(regexString string) (*ParseContext, error) {
	if len(regexString) == 0 {
		return nil, fmt.Errorf("missing regex string")
	}

	regex := []byte(regexString)
	ctx := &ParseContext{
		Pos:    0,
		Tokens: []token.Token{},
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
