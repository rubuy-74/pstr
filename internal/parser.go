package internal

import (
	"fmt"
	"slices"
	"strconv"
)

const infinite = -1

type tokenType uint8

const (
	group           tokenType = iota
	bracket         tokenType = iota
	or              tokenType = iota
	repeat          tokenType = iota
	literal         tokenType = iota
	groupUncaptured tokenType = iota
)

type token struct {
	tokenType tokenType
	value     any
}

type parseContext struct {
	pos    int
	tokens []token
}

func process(regex []byte, ctx *parseContext) error {
	ch := regex[ctx.pos]
	// fmt.Printf("character: %v\n", string(ch))
	// fmt.Printf("ctx.pos: %v\n", ctx.pos)
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
		processOr(regex, ctx)
	case '*':
		processRepeat(regex, ctx, 0, infinite)
	case '+':
		processRepeat(regex, ctx, 1, infinite)
	case '?':
		processRepeat(regex, ctx, 0, 1)
	case '{':
		minimum, maximum := getMinMaxRange(regex, ctx)
		fmt.Println(minimum, maximum)
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

func getMinMaxRange(regex []byte, ctx *parseContext) (minimum int, maximum int) {
	newPos := findNextSymbol(regex, ctx.pos, '}')
	rangeString := regex[ctx.pos+1 : newPos]

	// TODO: better checking (use macros)
	switch len(rangeString) {
	case 1:
		value, _ := strconv.Atoi(string(rangeString))
		return value, value
	case 2:
		minimum, _ := strconv.Atoi(string(rangeString[0]))
		return minimum, infinite
	case 3:
		minimum, _ := strconv.Atoi(string(rangeString[0]))
		maximum, _ := strconv.Atoi(string(rangeString[2]))
		return minimum, maximum
	}
	return infinite, infinite
}

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

func chunkBytes(data []byte, size int) [][]byte {
	numSlices := (len(data) + size - 1) / size
	subslices := make([][]byte, 0, numSlices)

	for i := 0; i < len(data); i += size {
		end := min(i+size, len(data))
		subslice := []byte{data[i], data[end-1]}
		subslices = append(subslices, subslice)
	}

	return subslices
}

func processGroup(regex []byte, ctx *parseContext) error {
	ctx.pos++
	newPos := findNextSymbol(regex, ctx.pos, ')')
	if newPos == 1 {
		return fmt.Errorf("invalid ( in the regex string")
	}
	groupRegex := regex[ctx.pos:newPos]
	groupCtx := &parseContext{
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

func processBrackets(regex []byte, ctx *parseContext) error {
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
			value: []byte{
				insideRegex[1],
				insideRegex[len(insideRegex)-2],
			},
		})
	}
	fmt.Println(ctx.tokens)

	ctx.pos = newPos
	return nil
}

func processOr(regex []byte, ctx *parseContext) {
	// TODO: implement this
}

func processRepeat(regex []byte, ctx *parseContext, min int, max int) {
	// TODO: implement this
}

func Parse(regexString string) (*parseContext, error) {
	regex := []byte(regexString)
	ctx := &parseContext{
		pos:    0,
		tokens: []token{},
	}
	for ctx.pos < len(regex) {
		err := process(regex, ctx)
		if err != nil {
			return nil, err
		}
		ctx.pos++
		// fmt.Println(ctx.tokens)
	}

	return ctx, nil
}
