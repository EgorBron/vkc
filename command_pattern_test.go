package vkc

import (
	"regexp"
	"testing"
)

func TestText(t *testing.T) {
	tests := []struct {
		name     string
		matcher  string
		input    string
		expected bool
	}{
		{
			name:     "exact match",
			matcher:  "help",
			input:    "help",
			expected: true,
		},
		{
			name:     "no match case sensitive",
			matcher:  "help",
			input:    "Help",
			expected: false,
		},
		{
			name:     "no match different text",
			matcher:  "help",
			input:    "hlep",
			expected: false,
		},
		{
			name:     "no match partial",
			matcher:  "help",
			input:    "help me",
			expected: false,
		},
		{
			name:     "empty strings",
			matcher:  "",
			input:    "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := Text(tt.matcher)
			result := pattern(tt.input)
			if result != tt.expected {
				t.Errorf("Text(%q)(%q) = %v, want %v", tt.matcher, tt.input, result, tt.expected)
			}
		})
	}
}

func TestListOf(t *testing.T) {
	tests := []struct {
		name     string
		matcher  []string
		input    string
		expected bool
	}{
		{
			name:     "match first",
			matcher:  []string{"help", "info", "about"},
			input:    "help",
			expected: true,
		},
		{
			name:     "match middle",
			matcher:  []string{"help", "info", "about"},
			input:    "info",
			expected: true,
		},
		{
			name:     "match last",
			matcher:  []string{"help", "info", "about"},
			input:    "about",
			expected: true,
		},
		{
			name:     "no match",
			matcher:  []string{"help", "info", "about"},
			input:    "unknown",
			expected: false,
		},
		{
			name:     "case sensitive",
			matcher:  []string{"help", "info"},
			input:    "Help",
			expected: false,
		},
		{
			name:     "empty list",
			matcher:  []string{},
			input:    "help",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := ListOf(tt.matcher)
			result := pattern(tt.input)
			if result != tt.expected {
				t.Errorf("ListOf(%v)(%q) = %v, want %v", tt.matcher, tt.input, result, tt.expected)
			}
		})
	}
}

func TestRegex(t *testing.T) {
	tests := []struct {
		name     string
		matcher  *regexp.Regexp
		input    string
		expected bool
	}{
		{
			name:     "simple match",
			matcher:  regexp.MustCompile(`^help$`),
			input:    "help",
			expected: true,
		},
		{
			name:     "digit pattern",
			matcher:  regexp.MustCompile(`^\d+$`),
			input:    "12345",
			expected: true,
		},
		{
			name:     "digit pattern no match",
			matcher:  regexp.MustCompile(`^\d+$`),
			input:    "abc",
			expected: false,
		},
		{
			name:     "word boundary",
			matcher:  regexp.MustCompile(`\bhelp\b`),
			input:    "help me",
			expected: true,
		},
		{
			name:     "word boundary no match",
			matcher:  regexp.MustCompile(`\bhelp\b`),
			input:    "helpful",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := Regex(tt.matcher)
			result := pattern(tt.input)
			if result != tt.expected {
				t.Errorf("Regex(%v)(%q) = %v, want %v", tt.matcher.String(), tt.input, result, tt.expected)
			}
		})
	}
}

func TestRegexCmd(t *testing.T) {
	tests := []struct {
		name     string
		matcher  string
		input    string
		expected bool
	}{
		{
			name:     "simple pattern",
			matcher:  `^help$`,
			input:    "help",
			expected: true,
		},
		{
			name:     "digit pattern",
			matcher:  `^\d+$`,
			input:    "12345",
			expected: true,
		},
		{
			name:     "digit pattern no match",
			matcher:  `^\d+$`,
			input:    "abc",
			expected: false,
		},
		{
			name:     "case insensitive",
			matcher:  `(?i)^help$`,
			input:    "Help",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := RegexStr(tt.matcher)
			result := pattern(tt.input)
			if result != tt.expected {
				t.Errorf("RegexCmd(%q)(%q) = %v, want %v", tt.matcher, tt.input, result, tt.expected)
			}
		})
	}
}
