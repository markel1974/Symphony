//go:build !windows

package asciirender

import (
	"os"
	"strings"
	"syscall"
	"unsafe"
)

type Termios struct {
	Iflag  uint64
	Oflag  uint64
	Cflag  uint64
	Lflag  uint64
	Cc     [20]uint8
	Ispeed uint64
	Ospeed uint64
}

func ioctl(fd uintptr, request uintptr, ifReq uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, request, ifReq)
	if errno != 0 {
		return errno
	}
	return nil
}

func MakeStdInRaw() error {
	const ioctlReadTerm = 0x40487413
	const ioctlWriteTerm = 0x80487414
	const IGNBRK = 0x1
	const BRKINT = 0x2
	const PARMRK = 0x8
	const ISTRIP = 0x20
	const INLCR = 0x40
	const IGNCR = 0x80
	const ICRNL = 0x100
	const IXON = 0x200
	const OPOST = 0x1
	const ECHO = 0x8
	const ECHONL = 0x10
	const ICANON = 0x100
	const ISIG = 0x80
	const IEXTEN = 0x400
	const CSIZE = 0x300
	const PARENB = 0x1000
	const CS8 = 0x300
	const VMIN = 0x10
	const VTIME = 0x11
	var term Termios
	fd := os.Stdin.Fd()
	if err := ioctl(fd, ioctlReadTerm, uintptr(unsafe.Pointer(&term))); err != nil {
		return err
	}
	term.Iflag &^= IGNBRK | BRKINT | PARMRK | ISTRIP | INLCR | IGNCR | ICRNL | IXON
	term.Oflag &^= OPOST
	term.Lflag &^= ECHO | ECHONL | ICANON | ISIG | IEXTEN
	term.Cflag &^= CSIZE | PARENB
	term.Cflag |= CS8
	term.Cc[VMIN] = 1
	term.Cc[VTIME] = 0

	if err := ioctl(fd, ioctlWriteTerm, uintptr(unsafe.Pointer(&term))); err != nil {
		return err
	}
	return nil
}

func detectSpecialTermColor(termVal string) (Level, bool) {
	if termVal == "" {
		return ColorLevelNone, false
	}
	if termVal == "screen" {
		return ColorLevelHundreds, false
	}
	if strings.Contains(termVal, "256color") {
		return ColorLevelHundreds, false
	}
	if strings.Contains(termVal, "xterm") {
		return ColorLevelHundreds, false
	}
	return ColorLevelBasic, false
}

func IsTerminal(fd uintptr) bool {
	return fd == uintptr(syscall.Stdout) || fd == uintptr(syscall.Stdin) || fd == uintptr(syscall.Stderr)
}
