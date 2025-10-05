package stub

import (
	"embed"
	"path/filepath"
	"strings"
)

// All is an embedded file system containing files from the "sources" and "tests" directories.
//
//go:embed sources/*.go
//go:embed tests/*.go
var All embed.FS

// FileInfo represents a file with its name and associated data in bytes.
type FileInfo struct {
	Name string
	Data []byte
}

// NewFileInfo creates a new FileInfo instance with the provided file name, reading its data from an embedded file system.
func NewFileInfo(name string) *FileInfo {
	data, _ := All.ReadFile(name)
	return &FileInfo{Name: name, Data: data}
}

// walk recursively traverses the directory starting from baseDir, collecting files with names prefixed by prefix into files.
// Returns an error if there is a problem reading the directory or files.
func walk(baseDir string, prefix string, files *[]string) error {
	data, err := All.ReadDir(baseDir)
	if err != nil {
		return err
	}
	for _, v := range data {
		target := filepath.Join(baseDir, v.Name())
		if v.IsDir() {
			if err = walk(target, prefix, files); err != nil {
				return err
			}
		}
		if !strings.HasPrefix(v.Name(), prefix) {
			continue
		}
		*(files) = append(*(files), target)
	}
	return nil
}
