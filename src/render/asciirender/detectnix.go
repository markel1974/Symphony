//go:build !windows

package asciirender

import (
	"strings"
	"syscall"
)

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
