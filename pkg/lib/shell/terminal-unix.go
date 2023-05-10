//go:build linux
// +build linux

package shellUtils

import (
	"syscall"
	"unsafe"
)

type window struct {
	Row int
	Col int
}

func WindowSize() window {
	win := window{0, 0}
	tio := syscall.TIOCGWINSZ

	res, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(tio),
		uintptr(unsafe.Pointer(&win)),
	)
	if int(res) == -1 {
		return window{0, 0}
	}

	return win
}
