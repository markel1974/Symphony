package bytecode

import (
	"fmt"
	"strings"
	"unicode"
)

// _controlCharNames is a map of control character runes to their escaped string representations.
var _controlCharNames = map[rune]string{
	0x00: "\\nul", // Null
	0x01: "\\soh", // Start of Heading
	0x02: "\\stx", // Start of Text
	0x03: "\\etx", // End of Text
	0x04: "\\eot", // End of Transmission
	0x05: "\\enq", // Enquiry
	0x06: "\\ack", // Acknowledge
	0x07: "\\bel", // Bell
	0x08: "\\bs",  // Backspace
	0x0B: "\\vt",  // Vertical Tab
	0x0C: "\\ff",  // Form Feed
	0x0E: "\\so",  // Shift Out
	0x0F: "\\si",  // Shift In
	0x10: "\\dle", // Data Link Escape
	0x11: "\\dc1", // Device Control 1
	0x12: "\\dc2", // Device Control 2
	0x13: "\\dc3", // Device Control 3
	0x14: "\\dc4", // Device Control 4
	0x15: "\\nak", // Negative Acknowledge
	0x16: "\\syn", // Synchronous Idle
	0x17: "\\etb", // End of Transmission Block
	0x18: "\\can", // Cancel
	0x19: "\\em",  // End of Medium
	0x1A: "\\sub", // Substitute
	0x1B: "\\esc", // Escape
	0x1C: "\\fs",  // File Separator
	0x1D: "\\gs",  // Group Separator
	0x1E: "\\rs",  // Record Separator
	0x1F: "\\us",  // Unit Separator
	0x7F: "\\del", // Delete
}

// EscapeNonPrintable converts non-printable characters in a string into their escaped representations.
func EscapeNonPrintable(s string) string {
	var builder strings.Builder
	for _, r := range s {
		switch r {
		case '\n':
			builder.WriteString("\\n")
		case '\t':
			builder.WriteString("\\t")
		case '\r':
			builder.WriteString("\\r")
		default:
			if name, ok := _controlCharNames[r]; ok {
				builder.WriteString(name)
			} else if unicode.IsPrint(r) {
				builder.WriteRune(r)
			} else {
				builder.WriteString(fmt.Sprintf("\\x%02x", r))
			}
		}
	}

	return builder.String()
}
