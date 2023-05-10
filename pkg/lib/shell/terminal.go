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
		uintptr(syscall.TIOCGWINSZ), //此参数,不同的操作系统可能不一样,例如:TIOCGWINSZ_OSX
		uintptr(unsafe.Pointer(&win)),
	)
	if int(res) == -1 {
		return window{0, 0}
	}

	return win
}
