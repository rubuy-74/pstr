// Made by Cursor Agent
package parser

import (
	"testing"

	tokenModel "github.com/rubuy-74/pstr/internal/models/token"
	"github.com/rubuy-74/pstr/internal/models/token_type"
)

// TestEmptyInputValidation tests the Parse() function with empty inputs
func TestEmptyInputValidation(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{"Empty string", "", true},
		{"Whitespace only", "   ", false}, // Should parse as literal spaces
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for input %q, but got none", tt.input)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for input %q: %v", tt.input, err)
			}
		})
	}
}

// TestProcessRepeatEdgeCases tests the processRepeat function with edge cases
func TestProcessRepeatEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		regex       string
		expectError bool
		description string
	}{
		{"Star with no preceding token", "*", true, "Should fail when * has no token to repeat"},
		{"Plus with no preceding token", "+", true, "Should fail when + has no token to repeat"},
		{"Question with no preceding token", "?", true, "Should fail when ? has no token to repeat"},
		{"Valid star", "a*", false, "Should work with valid preceding token"},
		{"Valid plus", "a+", false, "Should work with valid preceding token"},
		{"Valid question", "a?", false, "Should work with valid preceding token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.regex)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for %q (%s), but got none", tt.regex, tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %q (%s): %v", tt.regex, tt.description, err)
			}
		})
	}
}

// TestProcessOrEdgeCases tests the processOr function with edge cases
func TestProcessOrEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		regex       string
		expectError bool
		description string
	}{
		{"Or with no left operand", "|abc", true, "Should fail when | has no left operand"},
		{"Or with no right operand", "abc|", true, "Should fail when | has no right operand"},
		{"Or with both operands empty", "|", true, "Should fail when | has no operands"},
		{"Valid or", "a|b", false, "Should work with valid operands"},
		{"Or in group", "(a|b)", false, "Should work inside groups"},
		{"Complex or", "ab|cd", false, "Should work with complex operands"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.regex)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for %q (%s), but got none", tt.regex, tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %q (%s): %v", tt.regex, tt.description, err)
			}
		})
	}
}

// TestProcessBracketsEdgeCases tests the processBrackets function with edge cases
func TestProcessBracketsEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		regex       string
		expectError bool
		description string
	}{
		{"Empty brackets", "[]", true, "Should fail with empty bracket expression"},
		{"Single character brackets", "[a]", true, "Should fail with single character (needs at least 2)"},
		{"Valid range", "[a-z]", false, "Should work with valid range"},
		{"Unclosed brackets", "[a-z", true, "Should fail with unclosed brackets"},
		{"Invalid range", "[z-a]", false, "Should work even with reverse range"},
		{"Multiple ranges", "[a-zA-Z0-9]", false, "Should work with multiple ranges"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.regex)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for %q (%s), but got none", tt.regex, tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %q (%s): %v", tt.regex, tt.description, err)
			}
		})
	}
}

// TestProcessGroupEdgeCases tests the processGroup function with edge cases
func TestProcessGroupEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		regex       string
		expectError bool
		description string
	}{
		{"Empty group", "()", true, "Should fail with empty group"},
		{"Unclosed group", "(abc", true, "Should fail with unclosed group"},
		{"Valid group", "(abc)", false, "Should work with valid group"},
		{"Nested groups", "((abc))", true, "Should fail with nested groups (not implemented)"},
		{"Group with or", "(a|b)", false, "Should work with or inside group"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.regex)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for %q (%s), but got none", tt.regex, tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %q (%s): %v", tt.regex, tt.description, err)
			}
		})
	}
}

// TestGetMinMaxRangeEdgeCases tests the getMinMaxRange function with edge cases
func TestGetMinMaxRangeEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		regex       string
		expectError bool
		description string
	}{
		{"Unclosed range", "a{2", true, "Should fail with unclosed range"},
		{"Empty range", "a{}", false, "Should work with empty range (treated as {0,0})"},
		{"Invalid range syntax", "a{2,3,4}", true, "Should fail with too many commas"},
		{"Valid fixed range", "a{2}", false, "Should work with fixed range"},
		{"Valid min range", "a{2,}", false, "Should work with min range"},
		{"Valid max range", "a{,3}", false, "Should work with max range"},
		{"Valid full range", "a{2,3}", false, "Should work with full range"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.regex)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for %q (%s), but got none", tt.regex, tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %q (%s): %v", tt.regex, tt.description, err)
			}
		})
	}
}

// TestParseResultValidation tests that Parse returns valid results for edge cases
func TestParseResultValidation(t *testing.T) {
	tests := []struct {
		name        string
		regex       string
		expectError bool
		description string
	}{
		{"Empty tokens", "", true, "Should fail when no tokens to create NFA"},
		{"Valid single token", "a", false, "Should work with single token"},
		{"Valid multiple tokens", "abc", false, "Should work with multiple tokens"},
		{"Valid with repetition", "a*", false, "Should work with repetition"},
		{"Valid with brackets", "[a-z]", false, "Should work with brackets"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := Parse(tt.regex)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for %q (%s), but got none", tt.regex, tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %q (%s): %v", tt.regex, tt.description, err)
			}
			if !tt.expectError && ctx == nil {
				t.Errorf("Expected ParseContext for %q (%s), but got nil", tt.regex, tt.description)
			}
		})
	}
}

// TestTokenToNFAEdgeCases tests the Token.ToNFA function with edge cases
func TestTokenToNFAEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		token       tokenModel.Token
		expectPanic bool
		description string
	}{
		{
			name: "Invalid Group Value",
			token: tokenModel.Token{
				TokenType: token_type.Group,
				Value:     "invalid", // Should be []Token
			},
			expectPanic: false, // Should not panic due to safe type assertion
			description: "Should handle invalid group value gracefully",
		},
		{
			name: "Invalid Bracket Value",
			token: tokenModel.Token{
				TokenType: token_type.Bracket,
				Value:     "invalid", // Should be []BracketPayload
			},
			expectPanic: false, // Should not panic due to safe type assertion
			description: "Should handle invalid bracket value gracefully",
		},
		{
			name: "Invalid Or Value",
			token: tokenModel.Token{
				TokenType: token_type.Or,
				Value:     "invalid", // Should be []Token
			},
			expectPanic: false, // Should not panic due to safe type assertion
			description: "Should handle invalid or value gracefully",
		},
		{
			name: "Invalid Repeat Value",
			token: tokenModel.Token{
				TokenType: token_type.Repeat,
				Value:     "invalid", // Should be RepeatPayload
			},
			expectPanic: false, // Should not panic due to safe type assertion
			description: "Should handle invalid repeat value gracefully",
		},
		{
			name: "Invalid Literal Value",
			token: tokenModel.Token{
				TokenType: token_type.Literal,
				Value:     "invalid", // Should be uint8
			},
			expectPanic: false, // Should not panic due to safe type assertion
			description: "Should handle invalid literal value gracefully",
		},
		{
			name: "Valid Literal",
			token: tokenModel.Token{
				TokenType: token_type.Literal,
				Value:     uint8('a'),
			},
			expectPanic: false,
			description: "Should work with valid literal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("Unexpected panic for %s: %v", tt.description, r)
					}
				} else if tt.expectPanic {
					t.Errorf("Expected panic for %s, but got none", tt.description)
				}
			}()

			start, end := tt.token.ToNFA()
			if start == nil || end == nil {
				t.Errorf("ToNFA returned nil states for %s", tt.description)
			}
		})
	}
}

// TestComplexEdgeCases tests complex combinations that previously caused crashes
func TestComplexEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		regex       string
		expectError bool
		description string
	}{
		{"Multiple operators without operands", "*+?", true, "Should fail with multiple operators"},
		{"Or with repetition", "a*|b+", false, "Should work with or and repetition"},
		{"Brackets with repetition", "[a-z]*", false, "Should work with brackets and repetition"},
		{"Group with repetition", "(ab)*", false, "Should work with group and repetition"},
		{"Nested brackets", "[[a-z]]", false, "Should work with nested brackets"},
		{"Complex valid regex", "(a|b)*[0-9]+", false, "Should work with complex valid regex"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.regex)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for %q (%s), but got none", tt.regex, tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %q (%s): %v", tt.regex, tt.description, err)
			}
		})
	}
}

// TestParseReliability tests the parser with edge cases that should not crash
func TestParseReliability(t *testing.T) {
	tests := []struct {
		name        string
		regex       string
		expectError bool
		description string
	}{
		{"Empty regex", "", true, "Should fail with empty regex"},
		{"Valid simple", "a", false, "Should work with simple literal"},
		{"Valid with repetition", "a*", false, "Should work with repetition"},
		{"Valid with brackets", "[a-z]", false, "Should work with brackets"},
		{"Valid with or", "a|b", false, "Should work with or"},
		{"Valid with group", "(ab)", false, "Should work with group"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := Parse(tt.regex)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for %q (%s), but got none", tt.regex, tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %q (%s): %v", tt.regex, tt.description, err)
			}
			if !tt.expectError && ctx == nil {
				t.Errorf("Expected ParseContext for %q (%s), but got nil", tt.regex, tt.description)
			}
		})
	}
}
