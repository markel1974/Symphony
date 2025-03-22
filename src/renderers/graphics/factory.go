package graphics

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/renderers/graphics/ascii_render"
	"github.com/markel1974/c64emu/src/renderers/graphics/gl_render"
	"strings"
)

// Factory is a type responsible for creating instances of IDisplayRender based on specified identifiers.
type Factory struct {
}

// NewFactory initializes and returns a pointer to a new instance of Factory.
func NewFactory() *Factory {
	f := &Factory{}
	return f
}

// Create generates an IDisplayRender implementation based on the provided `id`.
// Returns an ASCII-based renderer if `id` is "ascii", otherwise returns a GL-based renderer.
func (f *Factory) Create(id string) references.IDisplayRender {
	switch strings.ToLower(strings.TrimSpace(id)) {
	case "ascii":
		return ascii_render.New()
	case "gl":
		return gl_render.New()
	default:
		return gl_render.New()
	}
}
