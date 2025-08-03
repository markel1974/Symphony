package stats

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
	"math"
	"runtime"
)

// rtPlotData represents runtime plot data and configuration for rendering a dynamic graph.
// It includes data points, plot type, min/max values, and auto-scaling behavior.
type rtPlotData struct {
	rtPlotData   []float64
	rtPlotType   int
	rtPlotMinVal float64
	rtPlotMaxVal float64
	rtPlotAuto   bool
}

// CreateMemoryPlot creates a shell command for plotting real-time memory statistics including alloc, total, os, and GC data.
// It supports dynamic updates and allows controlling plot scaling or enabling auto-scaling via interactive inputs.
func CreateMemoryPlot() interfaces.ICommand {
	run := func(task interfaces.IProcess, args []string) error {
		plt := &rtPlotData{
			rtPlotType:   0,
			rtPlotAuto:   true,
			rtPlotData:   nil,
			rtPlotMinVal: math.Inf(1),
			rtPlotMaxVal: math.Inf(-1),
		}
		if len(args) > 0 {
			switch args[0] {
			case "alloc":
				plt.rtPlotType = 0
			case "total":
				plt.rtPlotType = 1
			case "os":
				plt.rtPlotType = 2
			case "gc":
				plt.rtPlotType = 3
			}
		}
		task.SetContext(plt)
		task.CreateTimer(0, 300, -1)
		return nil
	}
	readFn := func(task interfaces.IProcess, code int, key rune) {
		ctx := task.GetContext()
		plt := ctx.(*rtPlotData)

		interval := math.Abs(plt.rtPlotMaxVal - plt.rtPlotMinVal)
		scale := (interval * 10) / 100

		switch key {
		case 'a', '+':
			plt.rtPlotAuto = false
			plt.rtPlotMaxVal += scale
			plt.rtPlotMinVal -= scale
		case 'z', '-':
			plt.rtPlotAuto = false
			plt.rtPlotMaxVal -= scale
			plt.rtPlotMinVal += scale
		case 'r':
			plt.rtPlotAuto = !plt.rtPlotAuto
		}

	}
	timerFn := func(task interfaces.IProcess, tid int, interval int) {
		var m runtime.MemStats
		ctx := task.GetContext()
		plt := ctx.(*rtPlotData)

		runtime.ReadMemStats(&m)
		var val float64
		switch plt.rtPlotType {
		case 0:
			val = bToMb(m.Alloc)
		case 1:
			val = bToMb(m.TotalAlloc)
		case 2:
			val = bToMb(m.Sys)
		case 3:
			val = float64(m.NumGC)
		default:
			val = bToMb(m.Alloc)
		}

		if val < plt.rtPlotMinVal {
			plt.rtPlotMinVal = val
		}
		if val > plt.rtPlotMaxVal {
			plt.rtPlotMaxVal = val
		}

		plt.rtPlotData = append(plt.rtPlotData, val)
		if len(plt.rtPlotData) > 10 {
			plt.rtPlotData = plt.rtPlotData[1:]
		}

		task.PaintRequest()
	}
	paintFn := func(task interfaces.IProcess, surface interfaces.ISurface) {
		var minPlot float64 = 0
		var maxPlot float64 = 0
		ctx := task.GetContext()
		plt := ctx.(*rtPlotData)
		if !plt.rtPlotAuto {
			minPlot = plt.rtPlotMinVal
			maxPlot = plt.rtPlotMaxVal
		}
		surface.DrawSeries(plt.rtPlotData, -1, -1, minPlot, maxPlot)
	}
	root := process.NewCommand("rtplot", interfaces.CommandTypeFile, nil, true, run)
	root.SetHelp("Runtime Plot", "Runtime Plot")
	root.SetTimerFn(timerFn)
	root.SetPaintFn(paintFn)
	root.SetReadFn(readFn)

	return root
}
