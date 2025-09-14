// Made by Cursor Agent
package state_machine

import (
	"testing"

	"github.com/rubuy-74/pstr/internal/models/token"
	"github.com/rubuy-74/pstr/internal/models/token_type"
	"github.com/rubuy-74/pstr/internal/parser"
)

// TestToNFAEmptyTokens tests the ToNFA function with empty token lists
func TestToNFAEmptyTokens(t *testing.T) {
	tests := []struct {
		name        string
		ctx         *parser.ParseContext
		expectError bool
		description string
	}{
		{
			name: "Empty tokens",
			ctx: &parser.ParseContext{
				Pos:    0,
				Tokens: []token.Token{},
			},
			expectError: true,
			description: "Should fail when no tokens exist",
		},
		{
			name: "Nil tokens",
			ctx: &parser.ParseContext{
				Pos:    0,
				Tokens: nil,
			},
			expectError: true,
			description: "Should fail when tokens is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ToNFA(tt.ctx)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for %s, but got none", tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.description, err)
			}
		})
	}
}

// TestToNFAValidTokens tests the ToNFA function with valid token lists
func TestToNFAValidTokens(t *testing.T) {
	tests := []struct {
		name        string
		ctx         *parser.ParseContext
		expectError bool
		description string
	}{
		{
			name: "Single token",
			ctx: &parser.ParseContext{
				Pos: 0,
				Tokens: []token.Token{
					{
						TokenType: token_type.Literal,
						Value:     uint8('a'),
					},
				},
			},
			expectError: false,
			description: "Should work with single token",
		},
		{
			name: "Multiple tokens",
			ctx: &parser.ParseContext{
				Pos: 0,
				Tokens: []token.Token{
					{
						TokenType: token_type.Literal,
						Value:     uint8('a'),
					},
					{
						TokenType: token_type.Literal,
						Value:     uint8('b'),
					},
				},
			},
			expectError: false,
			description: "Should work with multiple tokens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nfa, err := ToNFA(tt.ctx)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for %s, but got none", tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.description, err)
			}
			if !tt.expectError && nfa == nil {
				t.Errorf("Expected NFA for %s, but got nil", tt.description)
			}
		})
	}
}
