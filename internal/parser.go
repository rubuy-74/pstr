package internal

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

	value any
}

type parseContext struct {
	pos    int
	tokens []token
}

func process(regex string, ctx *parseContext) {
	ch := regex[ctx.pos]
	switch ch {
	case '(':
		processGroup(regex, ctx)
	case '[':
		processBrackets(regex, ctx)
	case '|':
		processOr(regex, ctx)
	case '*':
		processRepeat(regex, ctx, 0, infinite)
	case '+':
		processRepeat(regex, ctx, 1, infinite)
	case '?':
		processRepeat(regex, ctx, 0, 1)
	case '{':
		processRepeatRange(regex, ctx)
	}
}

func processGroup(regex string, ctx *parseContext) {
	// TODO: implement this
}
func processBrackets(regex string, ctx *parseContext) {
	// TODO: implement this
}
func processOr(regex string, ctx *parseContext) {
	// TODO: implement this
}
func processRepeat(regex string, ctx *parseContext, min int, max int) {
	// TODO: implement this
}
func processRepeatRange(regex string, ctx *parseContext) {
	// TODO: implement this
}

func Parse(regex string) *parseContext {
	ctx := &parseContext{
		pos:    0,
		tokens: []token{},
	}
	for ctx.pos < len(regex) {
		process(regex, ctx)
		ctx.pos++
	}

	return ctx
}
