package parser

import "testing"

func TestParseLiteral(t *testing.T) {
	ctx, err := Parse("abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.tokens) != 3 {
		t.Errorf("expected 3 tokens, got %d", len(ctx.tokens))
	}
}

func TestParseGroup(t *testing.T) {
	ctx, err := Parse("(ab)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.tokens) != 2 {
		t.Errorf("expected 2 tokens inside group, got %d", len(ctx.tokens))
	}
}

func TestParseBrackets(t *testing.T) {
	ctx, err := Parse("[a-z]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.tokens) == 0 || ctx.tokens[0].tokenType != bracket {
		t.Errorf("expected a bracket token, got %+v", ctx.tokens)
	}
}

func TestParseOr(t *testing.T) {
	ctx, err := Parse("a|b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ctx.tokens) != 1 {
		t.Errorf("expected one token, got %+v", ctx.tokens)
	}

	tokenOr := ctx.tokens[0]
	exprs, ok := tokenOr.value.([]token)
	if !ok {
		t.Errorf("expected token slice for Or operation's children, got %+v", tokenOr.value)
	}
	if len(exprs) != 2 {
		t.Errorf("expected two expressions, got %+v", ctx.tokens)
	}
	left := exprs[0]
	right := exprs[1]
	if left.tokenType != groupUncaptured || right.tokenType != groupUncaptured {
		t.Errorf("expected expressions to be valid, got left:{{ %+v }} right:{{ %+v }}", left, right)
	}
}

func TestParseOrComplex(t *testing.T) {
	ctx, err := Parse("b[a-zA-Z]|c{2,}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ctx.tokens) != 1 {
		t.Errorf("expected one token, got %+v", ctx.tokens)
	}

	tokenOr := ctx.tokens[0]
	exprs, ok := tokenOr.value.([]token)
	if !ok {
		t.Errorf("expected token slice for Or operation's children, got %+v", tokenOr.value)
	}
	if len(exprs) != 2 {
		t.Errorf("expected two expressions, got %+v", ctx.tokens)
	}
	left := exprs[0]
	right := exprs[1]
	if left.tokenType != groupUncaptured || right.tokenType != groupUncaptured {
		t.Errorf("expected expressions to be valid, got left:{{ %+v }} right:{{ %+v }}", left, right)
	}
}

func TestParseRepetitionStar(t *testing.T) {
	ctx, err := Parse("a*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.tokens) != 1 {
		t.Errorf("expected one token, got %+v", ctx.tokens)
	}

	token := ctx.tokens[0]
	tokenType := token.tokenType

	if tokenType != repeat {
		t.Errorf("expected token type = repeat, got %+v", tokenType)
	}

	tokenValue, ok := token.value.(repeatPayload)
	if !ok {
		t.Errorf("invalid token value, got %+v", token.value)
	}
	if tokenValue.min != 0 || tokenValue.max != -1 {
		t.Errorf("unexpected min max values, got min:%+v max: %+v", tokenValue.min, tokenValue.max)
	}
}

func TestParseRepetitionPlus(t *testing.T) {
	ctx, err := Parse("a+")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.tokens) != 1 {
		t.Errorf("expected one token, got %+v", ctx.tokens)
	}

	token := ctx.tokens[0]
	tokenType := token.tokenType

	if tokenType != repeat {
		t.Errorf("expected token type = repeat, got %+v", tokenType)
	}

	tokenValue, ok := token.value.(repeatPayload)
	if !ok {
		t.Errorf("invalid token value, got %+v", token.value)
	}
	if tokenValue.min != 1 || tokenValue.max != -1 {
		t.Errorf("unexpected min max values, got min:%+v max: %+v", tokenValue.min, tokenValue.max)
	}
}

func TestParseRepetitionSimpleRange(t *testing.T) {
	ctx, err := Parse("b{2}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.tokens) != 1 {
		t.Errorf("expected one token, got %+v", ctx.tokens)
	}

	token := ctx.tokens[0]
	tokenType := token.tokenType

	if tokenType != repeat {
		t.Errorf("expected token type = repeat, got %+v", tokenType)
	}

	tokenValue, ok := token.value.(repeatPayload)
	if !ok {
		t.Errorf("invalid token value, got %+v", token.value)
	}
	if tokenValue.min != 2 || tokenValue.max != 2 {
		t.Errorf("unexpected min max values, got min:%+v max: %+v", tokenValue.min, tokenValue.max)
	}
}

func TestParseRepetitionMinInfinite(t *testing.T) {
	ctx, err := Parse("b{2,}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.tokens) != 1 {
		t.Errorf("expected one token, got %+v", ctx.tokens)
	}

	token := ctx.tokens[0]
	tokenType := token.tokenType

	if tokenType != repeat {
		t.Errorf("expected token type = repeat, got %+v", tokenType)
	}

	tokenValue, ok := token.value.(repeatPayload)
	if !ok {
		t.Errorf("invalid token value, got %+v", token.value)
	}
	if tokenValue.min != 2 || tokenValue.max != -1 {
		t.Errorf("unexpected min max values, got min:%+v max: %+v", tokenValue.min, tokenValue.max)
	}
}

func TestParseRepetitionInfiniteMax(t *testing.T) {
	ctx, err := Parse("b{,2}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.tokens) != 1 {
		t.Errorf("expected one token, got %+v", ctx.tokens)
	}

	token := ctx.tokens[0]
	tokenType := token.tokenType

	if tokenType != repeat {
		t.Errorf("expected token type = repeat, got %+v", tokenType)
	}

	tokenValue, ok := token.value.(repeatPayload)
	if !ok {
		t.Errorf("invalid token value, got %+v", token.value)
	}
	if tokenValue.min != -1 || tokenValue.max != 2 {
		t.Errorf("unexpected min max values, got min:%+v max: %+v", tokenValue.min, tokenValue.max)
	}
}
