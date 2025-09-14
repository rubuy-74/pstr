package parser

import (
	"testing"

	tokenModel "github.com/rubuy-74/pstr/internal/models/token"
	"github.com/rubuy-74/pstr/internal/models/token_type"
)

func TestParseLiteral(t *testing.T) {
	ctx, err := Parse("abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.Tokens) != 3 {
		t.Errorf("expected 3 tokens, got %d", len(ctx.Tokens))
	}
}

func TestParseGroup(t *testing.T) {
	ctx, err := Parse("(ab)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.Tokens) != 2 {
		t.Errorf("expected 2 tokens inside group, got %d", len(ctx.Tokens))
	}
}

func TestParseBrackets(t *testing.T) {
	ctx, err := Parse("[a-z]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.Tokens) == 0 || ctx.Tokens[0].TokenType != token_type.Bracket {
		t.Errorf("expected a bracket token, got %+v", ctx.Tokens)
	}
}

func TestParseOr(t *testing.T) {
	ctx, err := Parse("a|b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ctx.Tokens) != 1 {
		t.Errorf("expected one token, got %+v", ctx.Tokens)
	}

	tokenOr := ctx.Tokens[0]
	exprs, ok := tokenOr.Value.([]tokenModel.Token)
	if !ok {
		t.Errorf("expected token slice for Or operation's children, got %+v", tokenOr.Value)
	}
	if len(exprs) != 2 {
		t.Errorf("expected two expressions, got %+v", ctx.Tokens)
	}
	left := exprs[0]
	right := exprs[1]
	if left.TokenType != token_type.GroupUncaptured || right.TokenType != token_type.GroupUncaptured {
		t.Errorf("expected expressions to be valid, got left:{{ %+v }} right:{{ %+v }}", left, right)
	}
}

func TestParseOrComplex(t *testing.T) {
	ctx, err := Parse("b[a-zA-Z]|c{2,}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ctx.Tokens) != 1 {
		t.Errorf("expected one token, got %+v", ctx.Tokens)
	}

	tokenOr := ctx.Tokens[0]
	exprs, ok := tokenOr.Value.([]tokenModel.Token)
	if !ok {
		t.Errorf("expected token slice for Or operation's children, got %+v", tokenOr.Value)
	}
	if len(exprs) != 2 {
		t.Errorf("expected two expressions, got %+v", ctx.Tokens)
	}
	left := exprs[0]
	right := exprs[1]
	if left.TokenType != token_type.GroupUncaptured || right.TokenType != token_type.GroupUncaptured {
		t.Errorf("expected expressions to be valid, got left:{{ %+v }} right:{{ %+v }}", left, right)
	}
}

func TestParseRepetitionStar(t *testing.T) {
	ctx, err := Parse("a*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.Tokens) != 1 {
		t.Errorf("expected one token, got %+v", ctx.Tokens)
	}

	token := ctx.Tokens[0]
	tokenType := token.TokenType

	if tokenType != token_type.Repeat {
		t.Errorf("expected token type = token_type.Repeat, got %+v", tokenType)
	}

	tokenValue, ok := token.Value.(tokenModel.RepeatPayload)
	if !ok {
		t.Errorf("invalid token.Value, got %+v", token.Value)
	}
	if tokenValue.Min != 0 || tokenValue.Max != -1 {
		t.Errorf("unexpected min max values, got min:%+v max: %+v", tokenValue.Min, tokenValue.Max)
	}
}

func TestParseRepetitionPlus(t *testing.T) {
	ctx, err := Parse("a+")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.Tokens) != 1 {
		t.Errorf("expected one token, got %+v", ctx.Tokens)
	}

	token := ctx.Tokens[0]
	tokenType := token.TokenType

	if tokenType != token_type.Repeat {
		t.Errorf("expected token type = token_type.Repeat, got %+v", tokenType)
	}

	tokenValue, ok := token.Value.(tokenModel.RepeatPayload)
	if !ok {
		t.Errorf("invalid token.Value, got %+v", token.Value)
	}
	if tokenValue.Min != 1 || tokenValue.Max != -1 {
		t.Errorf("unexpected min max values, got min:%+v max: %+v", tokenValue.Min, tokenValue.Max)
	}
}

func TestParseRepetitionSimpleRange(t *testing.T) {
	ctx, err := Parse("b{2}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.Tokens) != 1 {
		t.Errorf("expected one token, got %+v", ctx.Tokens)
	}

	token := ctx.Tokens[0]
	tokenType := token.TokenType

	if tokenType != token_type.Repeat {
		t.Errorf("expected token type = token_type.Repeat, got %+v", tokenType)
	}

	tokenValue, ok := token.Value.(tokenModel.RepeatPayload)
	if !ok {
		t.Errorf("invalid token.Value, got %+v", token.Value)
	}
	if tokenValue.Min != 2 || tokenValue.Max != 2 {
		t.Errorf("unexpected min max values, got min:%+v max: %+v", tokenValue.Min, tokenValue.Max)
	}
}

func TestParseRepetitionMinInfinite(t *testing.T) {
	ctx, err := Parse("b{2,}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.Tokens) != 1 {
		t.Errorf("expected one token, got %+v", ctx.Tokens)
	}

	token := ctx.Tokens[0]
	tokenType := token.TokenType

	if tokenType != token_type.Repeat {
		t.Errorf("expected token type = token_type.Repeat, got %+v", tokenType)
	}

	tokenValue, ok := token.Value.(tokenModel.RepeatPayload)
	if !ok {
		t.Errorf("invalid token.Value, got %+v", token.Value)
	}
	if tokenValue.Min != 2 || tokenValue.Max != -1 {
		t.Errorf("unexpected min max values, got min:%+v max: %+v", tokenValue.Min, tokenValue.Max)
	}
}

func TestParseRepetitionInfiniteMax(t *testing.T) {
	ctx, err := Parse("b{,2}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.Tokens) != 1 {
		t.Errorf("expected one token, got %+v", ctx.Tokens)
	}

	token := ctx.Tokens[0]
	tokenType := token.TokenType

	if tokenType != token_type.Repeat {
		t.Errorf("expected token type = token_type.Repeat, got %+v", tokenType)
	}

	tokenValue, ok := token.Value.(tokenModel.RepeatPayload)
	if !ok {
		t.Errorf("invalid token.Value, got %+v", token.Value)
	}
	if tokenValue.Min != -1 || tokenValue.Max != 2 {
		t.Errorf("unexpected min max values, got min:%+v max: %+v", tokenValue.Min, tokenValue.Max)
	}
}
