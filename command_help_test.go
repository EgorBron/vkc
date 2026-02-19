package vkc

import (
	"testing"
)

func TestCommandHelp(t *testing.T) {
	tests := []struct {
		name    string
		help    CommandHelp
		title   string
		brief   string
		usage   string
		aliases string
		hidden  bool
	}{
		{
			name: "complete help",
			help: CommandHelp{
				Title:   "Help Command",
				Brief:   "Show help information",
				Usage:   "help [command]",
				Aliases: "h, ?",
				Hidden:  false,
			},
			title:   "Help Command",
			brief:   "Show help information",
			usage:   "help [command]",
			aliases: "h, ?",
			hidden:  false,
		},
		{
			name: "hidden command",
			help: CommandHelp{
				Title:  "Secret",
				Brief:  "Secret command",
				Hidden: true,
			},
			title:  "Secret",
			brief:  "Secret command",
			hidden: true,
		},
		{
			name: "empty help",
			help: CommandHelp{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.help.Title != tt.title {
				t.Errorf("Title = %q, want %q", tt.help.Title, tt.title)
			}
			if tt.help.Brief != tt.brief {
				t.Errorf("Brief = %q, want %q", tt.help.Brief, tt.brief)
			}
			if tt.help.Usage != tt.usage {
				t.Errorf("Usage = %q, want %q", tt.help.Usage, tt.usage)
			}
			if tt.help.Aliases != tt.aliases {
				t.Errorf("Aliases = %q, want %q", tt.help.Aliases, tt.aliases)
			}
			if tt.help.Hidden != tt.hidden {
				t.Errorf("Hidden = %v, want %v", tt.help.Hidden, tt.hidden)
			}
		})
	}
}
