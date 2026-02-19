package vkc

import (
	"testing"
)

func TestSplitArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single argument",
			input:    "arg1",
			expected: []string{"arg1"},
		},
		{
			name:     "multiple arguments",
			input:    "arg1 arg2 arg3",
			expected: []string{"arg1", "arg2", "arg3"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: []string{},
		},
		{
			name:     "leading spaces",
			input:    "   arg1 arg2",
			expected: []string{"arg1", "arg2"},
		},
		{
			name:     "trailing spaces",
			input:    "arg1 arg2   ",
			expected: []string{"arg1", "arg2"},
		},
		{
			name:     "multiple spaces between args",
			input:    "arg1   arg2    arg3",
			expected: []string{"arg1", "arg2", "arg3"},
		},
		{
			name:     "tabs and spaces",
			input:    "arg1\targ2  arg3",
			expected: []string{"arg1", "arg2", "arg3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitArgs(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("SplitArgs(%q) length = %d, want %d", tt.input, len(result), len(tt.expected))
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("SplitArgs(%q)[%d] = %q, want %q", tt.input, i, v, tt.expected[i])
				}
			}
		})
	}
}
