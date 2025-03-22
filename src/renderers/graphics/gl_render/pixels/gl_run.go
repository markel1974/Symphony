package pixels

import (
	"fmt"
	"github.com/markel1974/c64emu/src/renderers/graphics/gl_render/pixels/executor"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// GLRun initializes GLFW, sets up the graphics thread, and runs the provided main function in the graphics thread context.
// Returns an error if GLFW initialization fails or the graphics thread encounters an error during execution.
func GLRun(main func()) error {
	err := glfw.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize glfw: %v", err)
	}
	defer glfw.Terminate()
	return executor.GraphicThread.Run(main)
}
