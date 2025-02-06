package pixels

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/markel1974/c64emu/src/pixels/executor"
)

func GLRun(main func()) error {
	err := glfw.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize glfw: %v", err)
	}
	defer glfw.Terminate()
	return executor.GraphicThread.Run(main)
}
