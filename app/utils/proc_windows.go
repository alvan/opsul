package utils

import (
	"os"
	"syscall"
)

func ProcAttrGroup() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}

func ProcKill(pid int) error {
	proc, err := os.FindProcess(-pid)
	if err == nil {
		err = proc.Signal(syscall.SIGTERM)
	}

	return err
}
