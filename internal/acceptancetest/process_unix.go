//go:build unix

package acceptancetest

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
)

func configureAcceptanceProcess(command *exec.Cmd) {
	// Give the worker its own process group, inherited by the test binary and Terraform.
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func stopAcceptanceProcess(command *exec.Cmd) error {
	if command.Process == nil {
		return nil
	}
	// A negative PID targets the process group instead of only its root process.
	err := syscall.Kill(-command.Process.Pid, syscall.SIGKILL)
	if errors.Is(err, syscall.ESRCH) {
		return os.ErrProcessDone
	}
	return err
}
