package adapters

import (
	"errors"
	"fmt"
)

// StOk indicates no error.
// StReadTimeout indicates a timeout occurred during reading.
// StTimeout indicates a general timeout condition.
// StEof indicates the end of file has been reached.
// StNotPresent indicates the device is not present.
const (
	StOk          uint8 = 0    // No error
	StReadTimeout       = 0x02 // Timeout on reading
	StTimeout           = 0x03 // Timeout
	StEof               = 0x40 // End of file
	StNotPresent        = 0x80 // Device not present
)

type FType uint8

// FTypeDel represents a deleted file type.
// FTypeSeq represents a sequential file type.
// FTypePrg represents a program file type.
// FTypeUsr represents a user file type.
// FTypeRel represents a relative file type.
// FTypeUnk represents an unknown file type.
const (
	FTypeUnk FType = iota
	FTypeDel       // Deleted
	FTypeSeq       // Sequential
	FTypePrg       // Program
	FTypeUsr       // User
	FTypeRel       // Relative
)

type FMode uint8

// FModeRead represents the file mode for reading.
// FModeWrite represents the file mode for writing.
// FModeAppend represents the file mode for appending.
// FModeM represents the file mode for reading an open file.
const (
	FModeUnknown FMode = iota
	FModeRead          // Read
	FModeWrite         // Write
	FModeAppend        // Append
	FModeM             // Read open file
)

// ErrOk represents the error code for a successful operation.
// ErrScratched indicates an error where files are scratched.
// ErrUnimplemented indicates that the operation is unimplemented.
// ErrRead20 represents an error for missing block headers.
// ErrRead21 indicates a read error due to a missing sync character.
// ErrRead22 signifies that a data block is not present during reading.
// ErrRead23 indicates a checksum error in a data block during read.
// ErrRead24 represents a byte decoding error during reading.
// ErrWrite25 signifies a write-verify error while writing data.
// ErrWriteProtect indicates that write protection is enabled.
// ErrRead27 represents a checksum error in a header during read.
// ErrWrite28 signifies an error due to a long data block during writing.
// ErrDiskId indicates a mismatch in disk IDs.
// ErrSyntax30 represents a general syntax error.
// ErrSyntax31 indicates an invalid command syntax error.
// ErrSyntax32 signifies a syntax error due to a command too long.
// ErrSyntax33 represents a syntax error with wildcards during writing.
// ErrSyntax34 indicates a syntax error for a missing file name.
// ErrWriteFileOpen signifies an error where a file is open during writing.
// ErrFileNotOpen indicates an error where a file needs to be open but isn't.
// ErrFileNotFound represents an error indicating the file is not found.
// ErrFileExists indicates an error where a file already exists.
// ErrFileType represents an error due to a file type mismatch.
// ErrNoBlock indicates that there is no block available or found.
// ErrIllegalTS represents an illegal track or sector error.
// ErrNoChannel signifies that no channel is available or found.
// ErrDirError indicates a directory-related error.
// ErrDiskFull represents an error indicating that the disk is full.
// ErrStartup is the error for startup or a power-up message.
// ErrNotReady indicates an error when the drive is not ready.
const (
	ErrOk            = iota // 00 OK
	ErrScratched            // 01 FILES SCRATCHED
	ErrUnimplemented        // 03 UNIMPLEMENTED
	ErrRead20               // 20 READ ERROR (block header not found)
	ErrRead21               // 21 READ ERROR (no sync character)
	ErrRead22               // 22 READ ERROR (data block not present)
	ErrRead23               // 23 READ ERROR (checksum error in data block)
	ErrRead24               // 24 READ ERROR (byte decoding error)
	ErrWrite25              // 25 WRITE ERROR (write-verify error)
	ErrWriteProtect         // 26 WRITE PROTECT ON
	ErrRead27               // 27 READ ERROR (checksum error in header)
	ErrWrite28              // 28 WRITE ERROR (long data block)
	ErrDiskId               // 29 DISK ID MISMATCH
	ErrSyntax30             // 30 SYNTAX ERROR (general syntax)
	ErrSyntax31             // 31 SYNTAX ERROR (invalid command)
	ErrSyntax32             // 32 SYNTAX ERROR (command too long)
	ErrSyntax33             // 33 SYNTAX ERROR (wildcards on writing)
	ErrSyntax34             // 34 SYNTAX ERROR (missing file name)
	ErrWriteFileOpen        // 60 WRITE FILE OPEN
	ErrFileNotOpen          // 61 FILE NOT OPEN
	ErrFileNotFound         // 62 FILE NOT FOUND
	ErrFileExists           // 63 FILE EXISTS
	ErrFileType             // 64 FILE TYPE MISMATCH
	ErrNoBlock              // 65 NO BLOCK
	ErrIllegalTS            // 66 ILLEGAL TRACK OR SECTOR
	ErrNoChannel            // 70 NO CHANNEL
	ErrDirError             // 71 DIR ERROR
	ErrDiskFull             // 72 DISK FULL
	ErrStartup              // 73 Power-up message
	ErrNotReady             // 74 DRIVE NOT READY
)

// _baseErrors is a mapping of error codes to their corresponding error messages used to describe file or disk errors.
var _baseErrors = map[int]string{
	ErrOk:            "OK",
	ErrScratched:     "FILES SCRATCHED",
	ErrUnimplemented: "UNIMPLEMENTED",
	ErrRead20:        "READ ERROR",
	ErrRead21:        "READ ERROR",
	ErrRead22:        "READ ERROR",
	ErrRead23:        "READ ERROR",
	ErrRead24:        "READ ERROR",
	ErrWrite25:       "WRITE ERROR",
	ErrWriteProtect:  "WRITE PROTECT",
	ErrRead27:        "READ ERROR",
	ErrWrite28:       "WRITE ERROR",
	ErrDiskId:        "DISK ID MISMATCH",
	ErrSyntax30:      "SYNTAX ERROR",
	ErrSyntax31:      "SYNTAX ERROR",
	ErrSyntax32:      "SYNTAX ERROR",
	ErrSyntax33:      "SYNTAX ERROR",
	ErrSyntax34:      "SYNTAX ERROR",
	ErrWriteFileOpen: "WRITE FILE OPEN",
	ErrFileNotOpen:   "FILE NOT OPEN",
	ErrFileNotFound:  "FILE NOT FOUND",
	ErrFileExists:    "FILE EXISTS",
	ErrFileType:      "FILE TYPE MISMATCH",
	ErrNoBlock:       "NO BLOCK",
	ErrIllegalTS:     "ILLEGAL TRACK OR SECTOR",
	ErrNoChannel:     "NO CHANNEL",
	ErrDirError:      "DIR ERROR",
	ErrDiskFull:      "DISK FULL",
	ErrStartup:       "CBM DOS V2.6 1541",
	ErrNotReady:      "DRIVE NOT READY",
}

// _errors is a map that associates integer error codes with their corresponding error objects.
var _errors map[int]error

// init initializes the internal error map by populating it with formatted errors derived from the _baseErrors map.
func init() {
	_errors = make(map[int]error)
	for idx, v := range _baseErrors {
		_errors[idx] = fmt.Errorf("%d, %s\x0d", idx, v)
	}
}

// Error retrieves an error from the _errors map by the given index or returns a default "INVALID INDEX" error if not found.
func Error(idx int) error {
	const invalidIndex = "1000, INVALID INDEX"
	err, ok := _errors[idx]
	if !ok {
		err = errors.New(invalidIndex)
	}
	return err
}
