//go:build !windows
// +build !windows

package shellUtils

import (
	"syscall"
	"unsafe"
)

type window struct {
	Row uint16
	Col uint16
}

func WindowSize() window {
	win := window{0, 0}

	res, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&win)),
	)
	if int(res) == -1 {
		return window{0, 0}
	}

	return win
}
