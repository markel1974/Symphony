package pixels

import (
	"fmt"
	"github.com/markel1974/c64emu/src/render/glrender/pixels/executor"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// GLRun initializes the GLFW library, handles its lifecycle, and runs the provided main function on the graphic thread.
// Returns an error if GLFW initialization fails or if running the main function encounters an issue.
func GLRun(main func()) error {
	err := glfw.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize glfw: %v", err)
	}
	defer glfw.Terminate()
	return executor.GraphicThread.Run(main)
}
