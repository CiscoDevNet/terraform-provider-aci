//go:build windows

package acceptancetest

import (
	"os/exec"
	"strconv"
	"syscall"
)

func configureAcceptanceProcess(command *exec.Cmd) {
	// Give the worker its own console process group; cancellation is handled below.
	command.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}

func stopAcceptanceProcess(command *exec.Cmd) error {
	if command.Process == nil {
		return nil
	}
	// /T includes descendant processes; /F forces termination of the worker's process tree.
	return exec.Command("taskkill", "/PID", strconv.Itoa(command.Process.Pid), "/T", "/F").Run()
}
