package adapters

import (
	"archive/zip"
	"bytes"
	"fmt"
	"hash/crc32"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"
)

// BufferFileInfo represents metadata about a file, typically used to describe files in an in-memory buffer or virtual file system.
// It implements standard file information methods such as Name, Size, Mode, ModTime, IsDir, and Sys.
// The type holds the file name, size in bytes, file mode, modification time, and whether the file represents a directory.
type BufferFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
}

// NewBufferFileInfo creates a new BufferFileInfo instance with the given name and data, initializing its metadata accordingly.
func NewBufferFileInfo(plainName string, data []byte) *BufferFileInfo {
	return &BufferFileInfo{
		name:    plainName,
		size:    int64(len(data)),
		mode:    0644,
		modTime: time.Now(),
		isDir:   false,
	}
}

func (fi *BufferFileInfo) Name() string       { return path.Base(fi.name) }
func (fi *BufferFileInfo) Size() int64        { return fi.size }
func (fi *BufferFileInfo) Mode() fs.FileMode  { return fi.mode }
func (fi *BufferFileInfo) ModTime() time.Time { return fi.modTime }
func (fi *BufferFileInfo) IsDir() bool        { return fi.isDir }
func (fi *BufferFileInfo) Sys() interface{}   { return nil }

// CreateZipHeader generates a zip.FileHeader for a file, specifying its name, data, and compression method.
// Returns the created zip.FileHeader and error if any validation or processing fails.
func CreateZipHeader(name string, data []byte, zipMode uint16) (*zip.FileHeader, error) {
	cleanedName := path.Clean(strings.TrimPrefix(name, "/"))
	if cleanedName == "." || strings.Contains(cleanedName, "..") {
		return nil, fmt.Errorf("invalid zip name: %s", name)
	}
	cleanedName = strings.ReplaceAll(cleanedName, "\\", "/")
	info := NewBufferFileInfo(cleanedName, data)
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return nil, err
	}
	header.Name = cleanedName //.Name()
	header.Method = zipMode
	header.CRC32 = crc32.Checksum(data, crc32.MakeTable(crc32.IEEE))
	if header.Method == zip.Store {
		header.CompressedSize64 = header.UncompressedSize64
	}
	header.Modified = info.ModTime()
	return header, nil
}

func ZipExtension() string { return ".zip" }

// Zip represents a ZIP archive file located at a specific path on the filesystem.
type Zip struct {
	path string
}

// NewZip initializes and returns a new Directory instance with the specified path. Returns an error if initialization fails.
func NewZip(path string) (*Zip, error) {
	z := &Zip{path: path}
	return z, nil
}

func (a *Zip) Extension() string { return ZipExtension() }

// Name returns the file path of the current Zip archive instance.
func (a *Zip) Name() string {
	return a.path
}

// ReadDir extracts and returns a list of file information from the zip archive.
func (a *Zip) ReadDir() ([]os.FileInfo, error) {
	zf, err := zip.OpenReader(a.path)
	if err != nil {
		return nil, err
	}
	defer zf.Close()
	var out []os.FileInfo
	for _, file := range zf.File {
		v := file.FileInfo()
		out = append(out, v)
	}
	return out, nil
}

// ReadFile extracts and reads the content of the specified file from the zip archive by its name.
func (a *Zip) ReadFile(plainName string) ([]byte, error) {
	zf, err := zip.OpenReader(a.path)
	if err != nil {
		return nil, err
	}
	defer zf.Close()
	var target *zip.File = nil
	for _, file := range zf.File {
		if string(CreateFileName(file.Name)) == plainName {
			target = file
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("file not found %s", plainName)
	}
	fc, err := target.Open()
	if err != nil {
		return nil, err
	}
	defer fc.Close()
	content, err := io.ReadAll(fc)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// WriteFile writes a new file with the specified name and content into the zip archive at the path defined in the Zip instance.
// It ensures that the provided data is compressed and added to the archive using the Deflate compression method.
// Returns an error if the zip file cannot be opened, modified, or written to, or if there are issues creating the file header.
func (a *Zip) WriteFile(plainName string, data []byte) error {
	zipReader, err := zip.OpenReader(a.path)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	zipFile, err := os.OpenFile(a.path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	header, err := CreateZipHeader(plainName, data, zip.Deflate)
	if err != nil {
		return fmt.Errorf("failed to create header from file info: %w", err)
	}
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	dataReader := bytes.NewReader(data)
	_, err = io.Copy(writer, dataReader)
	if err != nil {
		return err
	}
	return nil
}

// RenameFile renames a file from oldName to newName within the directory. Returns an error if the operation fails.
func (a *Zip) RenameFile(oldName string, newName string) error {
	return fmt.Errorf("unimplemented")
}

// Format applies a specific format operation to the directory, based on the provided name and id, and returns an error if unimplemented.
func (a *Zip) Format(name string, id string) error {
	return fmt.Errorf("unimplemented")
}

// Reset reinitializes the directory to a default or empty state, if supported. Returns an error if unimplemented.
func (a *Zip) Reset() error {
	return fmt.Errorf("unimplemented")
}

// Validate checks the state of the Directory instance and ensures it meets required conditions or constraints.
func (a *Zip) Validate() error {
	return fmt.Errorf("unimplemented")
}

// ScratchFile removes or marks a specified file for deletion within the directory, using the given command data.
func (a *Zip) ScratchFile(commandData string) error {
	return fmt.Errorf("unimplemented")
}

// BlockRead attempts to read a block of data from a specified track and sector but always returns an error for a Void adapter.
func (a *Zip) BlockRead(track int, sector int) error {
	return fmt.Errorf("unimplemented")
}

// BlockWrite writes a block to the specified track and sector and always returns an error for the Void adapter.
func (a *Zip) BlockWrite(track int, sector int) error {
	return fmt.Errorf("unimplemented")
}

// MemoryRead reads a specified memory region defined by the address and length. Returns an error for unsupported adapters.
func (a *Zip) MemoryRead(address uint16, length int) ([]byte, error) {
	return nil, fmt.Errorf("unimplemented")
}

// MemoryWrite writes a slice of bytes to the specified memory address and returns an error for void adapters.
func (a *Zip) MemoryWrite(address uint16, data []byte) error {
	return fmt.Errorf("unimplemented")
}

// Position sets the current position to the specified value, returning an error as the Void adapter is not functional.
func (a *Zip) Position(position int) error {
	return fmt.Errorf("unimplemented")
}
