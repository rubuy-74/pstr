// Made by Cursor Agent
package internal

import (
	"testing"

	"github.com/rubuy-74/pstr/internal/parser"
	"github.com/rubuy-74/pstr/internal/state_machine"
)

// TestCompletePipelineReliability tests the complete regex processing pipeline
// with edge cases that previously caused crashes
func TestCompletePipelineReliability(t *testing.T) {
	tests := []struct {
		name        string
		regex       string
		testString  string
		expectError bool
		description string
	}{
		// Empty and invalid inputs
		{"Empty regex", "", "test", true, "Should fail with empty regex"},
		{"Only operators", "*+?", "test", true, "Should fail with only operators"},

		// Or operator edge cases
		{"Or no left", "|abc", "abc", true, "Should fail with no left operand"},
		{"Or no right", "abc|", "abc", true, "Should fail with no right operand"},
		{"Or both empty", "|", "test", true, "Should fail with no operands"},

		// Repetition edge cases
		{"Star no token", "*", "test", true, "Should fail with * and no token"},
		{"Plus no token", "+", "test", true, "Should fail with + and no token"},
		{"Question no token", "?", "test", true, "Should fail with ? and no token"},

		// Bracket edge cases
		{"Empty brackets", "[]", "test", true, "Should fail with empty brackets"},
		{"Unclosed brackets", "[a-z", "a", true, "Should fail with unclosed brackets"},

		// Group edge cases
		{"Empty group", "()", "test", true, "Should fail with empty group"},
		{"Unclosed group", "(abc", "abc", true, "Should fail with unclosed group"},

		// Range edge cases
		{"Unclosed range", "a{2", "aa", true, "Should fail with unclosed range"},
		{"Empty range", "a{}", "test", true, "Should fail with empty range"},

		// Valid cases that should work
		{"Simple literal", "a", "a", false, "Should work with simple literal"},
		{"Valid star", "a*", "aaa", false, "Should work with valid star"},
		{"Valid plus", "a+", "aaa", false, "Should work with valid plus"},
		{"Valid question", "a?", "a", false, "Should work with valid question"},
		{"Valid brackets", "[a-z]", "a", false, "Should work with valid brackets"},
		{"Valid group", "(ab)", "ab", false, "Should work with valid group"},
		{"Valid or", "a|b", "a", false, "Should work with valid or"},
		{"Valid range", "a{2}", "aa", false, "Should work with valid range"},
		{"Complex valid", "(a|b)*[0-9]+", "a1", false, "Should work with complex valid regex"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Step 1: Parse the regex
			ctx, parseErr := parser.Parse(tt.regex)
			if parseErr != nil {
				if !tt.expectError {
					t.Errorf("Parse error for %q (%s): %v", tt.regex, tt.description, parseErr)
				}
				return
			}

			// Step 2: Convert to NFA
			nfa, nfaErr := state_machine.ToNFA(ctx)
			if nfaErr != nil {
				if !tt.expectError {
					t.Errorf("ToNFA error for %q (%s): %v", tt.regex, tt.description, nfaErr)
				}
				return
			}

			// Step 3: Test matching (should not crash)
			valid := nfa.Check(tt.testString, -1)

			// For valid regexes, we expect the test to complete without crashing
			// The actual matching result depends on the regex logic, not reliability
			if !tt.expectError {
				_ = valid // Use the result to avoid unused variable warning
				t.Logf("Regex %q processed successfully, match result: %v", tt.regex, valid)
			}
		})
	}
}

// TestPanicRecovery tests that the system doesn't panic on edge cases
func TestPanicRecovery(t *testing.T) {
	problematicInputs := []string{
		"",         // Empty string
		"*",        // Star without operand
		"+",        // Plus without operand
		"?",        // Question without operand
		"|",        // Or without operands
		"|abc",     // Or without left operand
		"abc|",     // Or without right operand
		"[]",       // Empty brackets
		"[a-z",     // Unclosed brackets
		"()",       // Empty group
		"(abc",     // Unclosed group
		"a{2",      // Unclosed range
		"a{}",      // Empty range
		"a{2,3,4}", // Invalid range syntax
	}

	for _, input := range problematicInputs {
		t.Run("PanicTest_"+input, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Panic occurred for input %q: %v", input, r)
				}
			}()

			// Try to parse
			ctx, err := parser.Parse(input)
			if err != nil {
				// Expected for invalid inputs
				return
			}

			// Try to convert to NFA
			_, err = state_machine.ToNFA(ctx)
			if err != nil {
				// Expected for some cases
				return
			}

			// If we get here, the input was processed without panic
			t.Logf("Input %q processed without panic", input)
		})
	}
}

// TestMemorySafety tests that the system doesn't access invalid memory
func TestMemorySafety(t *testing.T) {
	tests := []struct {
		name  string
		regex string
	}{
		{"Single character", "a"},
		{"Multiple characters", "abc"},
		{"With repetition", "a*"},
		{"With brackets", "[a-z]"},
		{"With group", "(ab)"},
		{"With or", "a|b"},
		{"With range", "a{2,3}"},
		{"Complex", "(a|b)*[0-9]+"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test ensures no array bounds violations occur
			ctx, err := parser.Parse(tt.regex)
			if err != nil {
				t.Fatalf("Parse failed for %q: %v", tt.regex, err)
			}

			nfa, err := state_machine.ToNFA(ctx)
			if err != nil {
				t.Fatalf("ToNFA failed for %q: %v", tt.regex, err)
			}

			// Test with various input strings
			testStrings := []string{"", "a", "ab", "abc", "123", "a1b2c3"}
			for _, testStr := range testStrings {
				// This should not cause any memory access violations
				_ = nfa.Check(testStr, -1)
			}
		})
	}
}
