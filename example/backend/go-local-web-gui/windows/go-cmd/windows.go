package cmd

import (
	"os"
	"os/exec"
	"syscall"
)

func terminateProcessWindows(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}

func setProcessGroupIDWindows(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{}
}
