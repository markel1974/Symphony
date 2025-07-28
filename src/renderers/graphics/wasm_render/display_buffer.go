package wasm_render

import "unsafe"

// DisplayBuffer is a structure that handles pixel data and optimized color representation for rendering a display surface.
// It includes buffers for individual RGBA colors and their expanded forms for multi-byte operations.
// The surface property maintains a byte slice representing the screen's current pixel data.
// maxLen determines the maximum allowable size of the surface buffer to prevent out-of-bounds operations.
type DisplayBuffer struct {
	w       int
	h       int
	surface []byte
}

// NewDisplayBuffer initializes and returns a pointer to a DisplayBuffer with specified width and height.
func NewDisplayBuffer(w int, h int) *DisplayBuffer {
	surface := make([]byte, w*h*4)
	return &DisplayBuffer{
		w:       w,
		h:       h,
		surface: surface,
	}
}

func (db *DisplayBuffer) GetWidth() int {
	return db.w
}

func (db *DisplayBuffer) GetHeight() int {
	return db.h
}

// GetSurface returns the internal surface buffer as a slice of bytes.
func (db *DisplayBuffer) GetSurface() []byte {
	return db.surface
}

// GetSurfaceLen returns the length of the DisplayBuffer's surface slice.
func (db *DisplayBuffer) GetSurfaceLen() int {
	return len(db.surface)
}

// GetSurfacePointer returns an unsafe.Pointer to the first byte of the surface slice in the DisplayBuffer.
func (db *DisplayBuffer) GetSurfacePointer() unsafe.Pointer {
	return unsafe.Pointer(&db.surface[0])
}

// SetArray updates a portion of the DisplayBuffer's surface starting at the specified index with data from the input slice.
// idx specifies the starting position in the surface buffer to begin the update.
// data points to the source slice containing byte data to copy into the surface.
// width determines the number of bytes to copy from the source slice into the surface buffer.
func (db *DisplayBuffer) SetArray(idx int, data *[]uint8, width int) {
	copy(db.surface[idx:], (*data)[:width])
}
