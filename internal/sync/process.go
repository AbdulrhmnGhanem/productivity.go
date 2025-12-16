package sync

import (
	"os"
	"os/exec"
	"syscall"
)

// TriggerBackgroundSync spawns a detached process to run 'readings sync'
func TriggerBackgroundSync() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.Command(exe, "sync")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true, // Detach from terminal
	}

	// We don't care about output
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	return cmd.Start()
}
