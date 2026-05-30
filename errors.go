package vkc

import (
	"fmt"
)

var (
	ErrCommandNotFound = fmt.Errorf("command not found")
	ErrNoPrefix        = fmt.Errorf("no prefix was found in message or the prefix matcher was not specified")
	ErrEmptyPrefix     = fmt.Errorf("prefix was empty")
	ErrNoPermissions   = fmt.Errorf("no permissions")
	ErrEmptyMessage    = fmt.Errorf("empty message")
)
