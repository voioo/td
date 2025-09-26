//go:build windows

package upgrade

import (
	"os/exec"
	"syscall"
)

// setWindowsHidden sets the command to run hidden on Windows
func setWindowsHidden(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
}
