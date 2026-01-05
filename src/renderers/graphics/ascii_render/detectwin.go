//go:build windows
// +build windows

package ascii_render

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	EnableVirtualTerminalProcessingMode uint32 = 0x4
)

var _kernel32 *syscall.LazyDLL
var _procGetConsoleMode *syscall.LazyProc
var _procSetConsoleMode *syscall.LazyProc
var _winVersion, _, _buildNumber = windows.RtlGetNtVersionNumbers()

var _isLikeInCmd bool
var _needVTP bool
var _enableWin bool

func init() {
	_isLikeInCmd = true
	_needVTP = false
	_enableWin = true
	if !SupportColor() {
		_isLikeInCmd = true
		return
	}
	if !_enableWin {
		return
	}
	tryEnableVTP(_needVTP)
}

func SupportColor() bool {
	return true
}

func tryEnableVTP(enable bool) bool {
	if !enable {
		return false
	}
	initKernel32Proc()
	if tryEnableOnCONOUT() {
		return true
	}
	return tryEnableOnStdout()
}

func initKernel32Proc() {
	if _kernel32 != nil {
		return
	}
	_kernel32 = syscall.NewLazyDLL("kernel32.dll")
	_procGetConsoleMode = _kernel32.NewProc("GetConsoleMode")
	_procSetConsoleMode = _kernel32.NewProc("SetConsoleMode")
}

func tryEnableOnCONOUT() bool {
	outHandle, err := syscall.Open("CONOUT$", syscall.O_RDWR, 0)
	if err != nil {
		return false
	}
	err = EnableVirtualTerminalProcessing(outHandle, true)
	if err != nil {
		return false
	}
	return true
}

func tryEnableOnStdout() bool {
	err := EnableVirtualTerminalProcessing(syscall.Stdout, true)
	if err != nil {
		return false
	}
	return true
}

func detectSpecialTermColor(termVal string) (tl Level, needVTP bool) {
	if os.Getenv("ConEmuANSI") == "ON" {
		return ColorLevelMillions, false
	}
	if _buildNumber < 10586 || _winVersion < 10 {
		if os.Getenv("ANSICON") != "" {
			conVersion := os.Getenv("ANSICON_VER")
			if conVersion >= "181" {
				return ColorLevelHundreds, false
			}
			return ColorLevelBasic, false
		}
		return ColorLevelNone, false
	}
	if _buildNumber < 14931 {
		return ColorLevelHundreds, true
	}
	return ColorLevelMillions, true
}

func EnableVirtualTerminalProcessing(stream syscall.Handle, enable bool) error {
	var mode uint32
	if err := syscall.GetConsoleMode(stream, &mode); err != nil {
		return err
	}
	if enable {
		mode |= EnableVirtualTerminalProcessingMode
	} else {
		mode &^= EnableVirtualTerminalProcessingMode
	}
	ret, _, err := procSetConsoleMode.Call(uintptr(stream), uintptr(mode))
	if ret == 0 {
		return err
	}
	return nil
}

func IsTty(fd uintptr) bool {
	initKernel32Proc()
	var st uint32
	r, _, e := syscall.Syscall(_procGetConsoleMode.Addr(), 2, fd, uintptr(unsafe.Pointer(&st)), 0)
	return r != 0 && e == 0
}

func IsTerminal(fd uintptr) bool {
	initKernel32Proc()
	var st uint32
	r, _, e := syscall.Syscall(_procGetConsoleMode.Addr(), 2, fd, uintptr(unsafe.Pointer(&st)), 0)
	return r != 0 && e == 0
}

func MakeStdInRaw() error {
	return EnableVirtualTerminalProcessing(syscall.Stdin, true)
}
