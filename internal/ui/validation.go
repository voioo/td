package ui

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// ValidateTaskName validates a task name input.
func ValidateTaskName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("task name cannot be empty")
	}
	if len(name) > 200 {
		return errors.New("task name is too long (maximum 200 characters)")
	}
	if !utf8.ValidString(name) {
		return errors.New("task name contains invalid characters")
	}
	// Check for potentially harmful characters
	if strings.ContainsAny(name, "\n\r\t") {
		return errors.New("task name cannot contain newlines or tabs")
	}
	return nil
}

// ValidatePriorityInput validates priority input.
func ValidatePriorityInput(priority int) error {
	if priority < 0 || priority > 3 {
		return errors.New("priority must be between 0 and 3")
	}
	return nil
}

// SanitizeTaskName sanitizes a task name by trimming whitespace and removing dangerous characters.
func SanitizeTaskName(name string) string {
	// Trim whitespace
	name = strings.TrimSpace(name)
	// Remove newlines and tabs
	name = strings.ReplaceAll(name, "\n", " ")
	name = strings.ReplaceAll(name, "\r", " ")
	name = strings.ReplaceAll(name, "\t", " ")
	// Collapse multiple spaces
	for strings.Contains(name, "  ") {
		name = strings.ReplaceAll(name, "  ", " ")
	}
	return strings.TrimSpace(name)
}
