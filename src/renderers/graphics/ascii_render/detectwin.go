//go:build windows
// +build windows

package ascii_render

import (
	"syscall"
	"unsafe"
)

const (
	EnableVirtualTerminalProcessingMode uint32 = 0x4
)

var kernel32 *syscall.LazyDLL
var procGetConsoleMode *syscall.LazyProc
var procSetConsoleMode *syscall.LazyProc
var winVersion, _, buildNumber = windows.RtlGetNtVersionNumbers()

func init() {
	if !SupportColor() {
		isLikeInCmd = true
		return
	}
	if !Enable {
		return
	}
	tryEnableVTP(needVTP)
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
	if kernel32 != nil {
		return
	}
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
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
	if buildNumber < 10586 || winVersion < 10 {
		if os.Getenv("ANSICON") != "" {
			conVersion := os.Getenv("ANSICON_VER")
			if conVersion >= "181" {
				return ColorLevelHundreds, false
			}
			return ColorLevelBasic, false
		}
		return ColorLevelNone, false
	}
	if buildNumber < 14931 {
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
	r, _, e := syscall.Syscall(procGetConsoleMode.Addr(), 2, fd, uintptr(unsafe.Pointer(&st)), 0)
	return r != 0 && e == 0
}

func IsTerminal(fd uintptr) bool {
	initKernel32Proc()
	var st uint32
	r, _, e := syscall.Syscall(procGetConsoleMode.Addr(), 2, fd, uintptr(unsafe.Pointer(&st)), 0)
	return r != 0 && e == 0
}
