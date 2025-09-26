package ui

import (
	"errors"
	"testing"
)

func TestValidateTaskName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected error
	}{
		{"valid name", "Buy groceries", nil},
		{"empty name", "", errTaskNameEmpty()},
		{"whitespace only", "   ", errTaskNameEmpty()},
		{"name too long", string(make([]byte, 201)), errTaskNameTooLong()},
		{"name with newline", "Task\nwith newline", errTaskNameInvalidChars()},
		{"name with tab", "Task\twith tab", errTaskNameInvalidChars()},
		{"valid unicode", "Comprar 🛒", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateTaskName(test.input)
			if test.expected == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if test.expected != nil && err == nil {
				t.Errorf("expected error %v, got nil", test.expected)
			}
			if test.expected != nil && err != nil && err.Error() != test.expected.Error() {
				t.Errorf("expected error %v, got %v", test.expected, err)
			}
		})
	}
}

func TestSanitizeTaskName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"trim spaces", "  task  ", "task"},
		{"remove newlines", "task\nwith\nlines", "task with lines"},
		{"remove tabs", "task\twith\ttabs", "task with tabs"},
		{"collapse spaces", "task  with   spaces", "task with spaces"},
		{"mixed cleanup", "  task\nwith\t  mixed  ", "task with mixed"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := SanitizeTaskName(test.input)
			if result != test.expected {
				t.Errorf("expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestValidatePriorityInput(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected error
	}{
		{"valid priority 0", 0, nil},
		{"valid priority 1", 1, nil},
		{"valid priority 2", 2, nil},
		{"valid priority 3", 3, nil},
		{"invalid negative", -1, errInvalidPriority()},
		{"invalid too high", 4, errInvalidPriority()},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidatePriorityInput(test.input)
			if test.expected == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if test.expected != nil && err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

// Helper functions to create expected errors
func errTaskNameEmpty() error {
	return errors.New("task name cannot be empty")
}

func errTaskNameTooLong() error {
	return errors.New("task name is too long (maximum 200 characters)")
}

func errTaskNameInvalidChars() error {
	return errors.New("task name cannot contain newlines or tabs")
}

func errInvalidPriority() error {
	return errors.New("priority must be between 0 and 3")
}
