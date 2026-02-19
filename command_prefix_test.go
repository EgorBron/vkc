package vkc

import (
	"regexp"
	"testing"
)

func TestPrefixText(t *testing.T) {
	tests := []struct {
		name           string
		matcher        string
		input          string
		expectedMatch  bool
		expectedRemain string
	}{
		{
			name:           "simple match",
			matcher:        "!",
			input:          "!help",
			expectedMatch:  true,
			expectedRemain: "help",
		},
		{
			name:           "match with spaces",
			matcher:        "!",
			input:          "!   help",
			expectedMatch:  true,
			expectedRemain: "help",
		},
		{
			name:           "no match",
			matcher:        "!",
			input:          "help",
			expectedMatch:  false,
			expectedRemain: "",
		},
		{
			name:           "prefix only",
			matcher:        "!",
			input:          "!",
			expectedMatch:  true,
			expectedRemain: "",
		},
		{
			name:           "multi-char prefix",
			matcher:        "prefix:",
			input:          "prefix: command arg",
			expectedMatch:  true,
			expectedRemain: "command arg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := PrefixText(tt.matcher)
			matched, remaining := matcher(tt.input)
			if matched != tt.expectedMatch {
				t.Errorf("PrefixText(%q)(%q) matched = %v, want %v", tt.matcher, tt.input, matched, tt.expectedMatch)
			}
			if remaining != tt.expectedRemain {
				t.Errorf("PrefixText(%q)(%q) remaining = %q, want %q", tt.matcher, tt.input, remaining, tt.expectedRemain)
			}
		})
	}
}

func TestPrefixListOf(t *testing.T) {
	tests := []struct {
		name           string
		matcher        []string
		input          string
		expectedMatch  bool
		expectedRemain string
	}{
		{
			name:           "match first prefix",
			matcher:        []string{"!", "/", "#"},
			input:          "!help me",
			expectedMatch:  true,
			expectedRemain: "help me",
		},
		{
			name:           "match middle prefix",
			matcher:        []string{"!", "/", "#"},
			input:          "/help",
			expectedMatch:  true,
			expectedRemain: "help",
		},
		{
			name:           "match last prefix",
			matcher:        []string{"!", "/", "#"},
			input:          "#help",
			expectedMatch:  true,
			expectedRemain: "help",
		},
		{
			name:           "no match",
			matcher:        []string{"!", "/", "#"},
			input:          "help",
			expectedMatch:  false,
			expectedRemain: "",
		},
		{
			name:           "with spaces",
			matcher:        []string{"!", "/"},
			input:          "!    help",
			expectedMatch:  true,
			expectedRemain: "help",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := PrefixListOf(tt.matcher)
			matched, remaining := matcher(tt.input)
			if matched != tt.expectedMatch {
				t.Errorf("PrefixListOf(%v)(%q) matched = %v, want %v", tt.matcher, tt.input, matched, tt.expectedMatch)
			}
			if remaining != tt.expectedRemain {
				t.Errorf("PrefixListOf(%v)(%q) remaining = %q, want %q", tt.matcher, tt.input, remaining, tt.expectedRemain)
			}
		})
	}
}

func TestPrefixRegex(t *testing.T) {
	tests := []struct {
		name           string
		matcher        *regexp.Regexp
		input          string
		expectedMatch  bool
		expectedRemain string
	}{
		{
			name:           "simple exclamation",
			matcher:        regexp.MustCompile(`^!(.*?)$`),
			input:          "!help me",
			expectedMatch:  true,
			expectedRemain: "help me",
		},
		{
			name:           "digit prefix",
			matcher:        regexp.MustCompile(`^\d+(.*?)$`),
			input:          "123help",
			expectedMatch:  true,
			expectedRemain: "help",
		},
		{
			name:           "no match",
			matcher:        regexp.MustCompile(`^!(.*?)$`),
			input:          "#help",
			expectedMatch:  false,
			expectedRemain: "",
		},
		{
			name:           "with spaces in capture",
			matcher:        regexp.MustCompile(`^!(.*?)$`),
			input:          "!   help me   ",
			expectedMatch:  true,
			expectedRemain: "help me",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := PrefixRegex(tt.matcher)
			matched, remaining := matcher(tt.input)
			if matched != tt.expectedMatch {
				t.Errorf("PrefixRegex(%v)(%q) matched = %v, want %v", tt.matcher.String(), tt.input, matched, tt.expectedMatch)
			}
			if remaining != tt.expectedRemain {
				t.Errorf("PrefixRegex(%v)(%q) remaining = %q, want %q", tt.matcher.String(), tt.input, remaining, tt.expectedRemain)
			}
		})
	}
}

func TestPrefixRegexCmd(t *testing.T) {
	tests := []struct {
		name           string
		matcher        string
		input          string
		expectedMatch  bool
		expectedRemain string
	}{
		{
			name:           "simple pattern",
			matcher:        `^!(.*?)$`,
			input:          "!help me",
			expectedMatch:  true,
			expectedRemain: "help me",
		},
		{
			name:           "digit prefix",
			matcher:        `^\d+(.*?)$`,
			input:          "123help",
			expectedMatch:  true,
			expectedRemain: "help",
		},
		{
			name:           "no match",
			matcher:        `^!(.*?)$`,
			input:          "#help",
			expectedMatch:  false,
			expectedRemain: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := PrefixRegexStr(tt.matcher)
			matched, remaining := matcher(tt.input)
			if matched != tt.expectedMatch {
				t.Errorf("PrefixRegexCmd(%q)(%q) matched = %v, want %v", tt.matcher, tt.input, matched, tt.expectedMatch)
			}
			if remaining != tt.expectedRemain {
				t.Errorf("PrefixRegexCmd(%q)(%q) remaining = %q, want %q", tt.matcher, tt.input, remaining, tt.expectedRemain)
			}
		})
	}
}

func TestPrefixFunc(t *testing.T) {
	tests := []struct {
		name           string
		matcher        func(string) (bool, string)
		input          string
		expectedMatch  bool
		expectedRemain string
	}{
		{
			name: "custom function",
			matcher: func(input string) (bool, string) {
				if len(input) > 0 && input[0] == '!' {
					return true, input[1:]
				}
				return false, ""
			},
			input:          "!help",
			expectedMatch:  true,
			expectedRemain: "help",
		},
		{
			name: "always fails",
			matcher: func(input string) (bool, string) {
				return false, ""
			},
			input:          "!help",
			expectedMatch:  false,
			expectedRemain: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := PrefixFunc(tt.matcher)
			matched, remaining := matcher(tt.input)
			if matched != tt.expectedMatch {
				t.Errorf("PrefixFunc - matched = %v, want %v", matched, tt.expectedMatch)
			}
			if remaining != tt.expectedRemain {
				t.Errorf("PrefixFunc - remaining = %q, want %q", remaining, tt.expectedRemain)
			}
		})
	}
}
