//go:build darwin
// +build darwin

package shellUtils

import (
	"runtime"
	"syscall"
	"unsafe"
)

type window struct {
	Row int
	Col int
}

func WindowSize() window {
	win := window{0, 0}
	tio := syscall.TIOCGWINSZ_OSX

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
