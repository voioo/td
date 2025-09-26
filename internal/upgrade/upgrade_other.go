//go:build !windows

package upgrade

import "os/exec"

// setWindowsHidden is a no-op on non-Windows platforms
func setWindowsHidden(cmd *exec.Cmd) {
	// No-op on non-Windows platforms
}
