package vkc

import (
	"testing"
)

func TestFindCommand(t *testing.T) {
	var nilexecutor = func(ctx CommandContext[any]) error { return nil }
	tests := []struct {
		name              string
		rawCmd            string
		handlers          []*CommandHandler[any]
		expectedHandler   *CommandHandler[any]
		expectedRemaining string
		shouldFindHandler bool
	}{
		{
			name:   "find single word command",
			rawCmd: "help",
			handlers: []*CommandHandler[any]{
				{
					Pattern:  Text("help"),
					Executor: nilexecutor,
				},
			},
			expectedRemaining: "",
			shouldFindHandler: true,
		},
		{
			name:   "find command with arguments",
			rawCmd: "help me something",
			handlers: []*CommandHandler[any]{
				{
					Pattern:  Text("help"),
					Executor: nilexecutor,
				},
			},
			expectedRemaining: "me something",
			shouldFindHandler: true,
		},
		{
			name:   "multi-word command",
			rawCmd: "user list admin",
			handlers: []*CommandHandler[any]{
				{
					Pattern:  Text("user list"),
					Executor: nilexecutor,
				},
			},
			expectedRemaining: "admin",
			shouldFindHandler: true,
		},
		{
			name:   "no matching command",
			rawCmd: "unknown command",
			handlers: []*CommandHandler[any]{
				{
					Pattern:  Text("help"),
					Executor: nilexecutor,
				},
			},
			expectedRemaining: "",
			shouldFindHandler: false,
		},
		{
			name:   "nil handler in list",
			rawCmd: "help",
			handlers: []*CommandHandler[any]{
				nil,
				{
					Pattern:  Text("help"),
					Executor: nilexecutor,
				},
			},
			expectedRemaining: "",
			shouldFindHandler: true,
		},
		{
			name:   "match first handler with longer command",
			rawCmd: "user list all",
			handlers: []*CommandHandler[any]{
				{
					Pattern:  Text("user list"),
					Executor: nilexecutor,
				},
				{
					Pattern:  Text("user"),
					Executor: nilexecutor,
				},
			},
			expectedRemaining: "all",
			shouldFindHandler: true,
		},
		{
			name:              "empty handlers list",
			rawCmd:            "help",
			handlers:          []*CommandHandler[any]{},
			expectedRemaining: "",
			shouldFindHandler: false,
		},
		{
			name:   "list pattern matcher",
			rawCmd: "info parameter",
			handlers: []*CommandHandler[any]{
				{
					Pattern:  ListOf([]string{"help", "info", "about"}),
					Executor: nilexecutor,
				},
			},
			expectedRemaining: "parameter",
			shouldFindHandler: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, remaining := FindCommand(tt.rawCmd, tt.handlers)

			if tt.shouldFindHandler {
				if handler == nil {
					t.Errorf("FindCommand(%q, ...) handler = nil, want non-nil", tt.rawCmd)
				}
				if remaining != tt.expectedRemaining {
					t.Errorf("FindCommand(%q, ...) remaining = %q, want %q", tt.rawCmd, remaining, tt.expectedRemaining)
				}
			} else {
				if handler != nil {
					t.Errorf("FindCommand(%q, ...) handler = %v, want nil", tt.rawCmd, handler)
				}
			}
		})
	}
}

func TestCommandHandlerIsAccessAvailable(t *testing.T) {
	tests := []struct {
		name           string
		accessCheck    *HandlerAccessCheck[any]
		expectedResult bool
	}{
		{
			name:           "no access check",
			accessCheck:    nil,
			expectedResult: true,
		},
		{
			name: "access check allows",
			accessCheck: &HandlerAccessCheck[any]{
				Checker: func(handler *CommandHandler[any], ctx CommandContext[any]) bool {
					return true
				},
			},
			expectedResult: true,
		},
		{
			name: "access check denies",
			accessCheck: &HandlerAccessCheck[any]{
				Checker: func(handler *CommandHandler[any], ctx CommandContext[any]) bool {
					return false
				},
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &CommandHandler[any]{
				Pattern:     Text("test"),
				Executor:    func(ctx CommandContext[any]) error { return nil },
				AccessCheck: tt.accessCheck,
			}

			ctx := CommandContext[any]{
				Dependency: nil,
			}

			result := handler.IsAccessAvailable(ctx)
			if result != tt.expectedResult {
				t.Errorf("IsAccessAvailable() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestCommandHandlerPatternMatchers(t *testing.T) {
	tests := []struct {
		name    string
		pattern CommandPattern
		input   string
		matches bool
	}{
		{
			name:    "text exact match",
			pattern: Text("command"),
			input:   "command",
			matches: true,
		},
		{
			name:    "text no match",
			pattern: Text("command"),
			input:   "other",
			matches: false,
		},
		{
			name:    "list match",
			pattern: ListOf([]string{"cmd1", "cmd2"}),
			input:   "cmd1",
			matches: true,
		},
		{
			name:    "list no match",
			pattern: ListOf([]string{"cmd1", "cmd2"}),
			input:   "cmd3",
			matches: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pattern(tt.input)
			if result != tt.matches {
				t.Errorf("Pattern(%q) = %v, want %v", tt.input, result, tt.matches)
			}
		})
	}
}
