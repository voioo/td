package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	// Save original args and restore them after test
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Test cases
	testCases := []struct {
		name string
		args []string
	}{
		{"version flag", []string{"td", "--version"}},
		{"version shorthand", []string{"td", "-v"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set test args
			os.Args = tc.args

			// We expect the program to exit with status 0
			// This is a simple test to ensure the flag is recognized
			if os.Getenv("TEST_VERSION_EXIT") == "1" {
				main()
				return
			}

			// Run the test in a subprocess
			executable, _ := os.Executable()
			cmd := exec.Command(executable, "-test.run=TestVersionFlag")
			cmd.Env = append(os.Environ(), "TEST_VERSION_EXIT=1")
			err := cmd.Run()

			if e, ok := err.(*exec.ExitError); ok && !e.Success() {
				t.Errorf("version flag test failed: %v", err)
			}
		})
	}
}
