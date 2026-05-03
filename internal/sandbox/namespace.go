//go:build linux

package sandbox

import (
	"os/exec"
	"syscall"
)

func buildSandboxedCmd(executable string, args []string) *exec.Cmd {
	cmd := exec.Command(executable, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNS,
		NoNewPrivs: true,
	}
	return cmd
}
