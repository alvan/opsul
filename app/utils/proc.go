// +build !windows

package utils

import "syscall"

func ProcAttrGroup() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}

func ProcKill(pid int) error {
	return syscall.Kill(-pid, syscall.SIGKILL)
}
