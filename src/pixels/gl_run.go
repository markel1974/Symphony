package pixels

import (
	"errors"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/markel1974/c64emu/src/pixels/executor"
)

func GLRun(run func()) {
	err := glfw.Init()
	if err != nil {
		panic(errors.New("failed to initialize glfw"))
	}
	defer glfw.Terminate()
	executor.GraphicThread.Run(run)
}

func GLGetTime() float64 {
	return glfw.GetTime()
}
