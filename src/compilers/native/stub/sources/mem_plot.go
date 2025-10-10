package sources

import (
	"fmt"
	"math"
	"runtime"
)

// ISurface represents an abstraction for rendering data-driven graphical series onto a surface.
// DrawSeries renders a series of data points within specified boundaries and dimensions.
type ISurface interface {
	DrawSeries(data []float64, rows int, columns int, min float64, max float64)
}

// bToMb converts a given size in bytes (b) to megabytes as a float64.
func bToMb(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

// MemPlot is a structure that tracks memory statistics over time for visualization and analysis.
type MemPlot struct {
	container []float64
	kind      int
	minVal    float64
	maxVal    float64
	auto      bool
}

// NewMemPlot initializes a new MemPlot instance of the specified kind, with auto scaling enabled and empty data.
func NewMemPlot(kind int) *MemPlot {
	plt := &MemPlot{
		kind:      kind,
		auto:      true,
		container: nil,
		minVal:    math.Inf(1),
		maxVal:    math.Inf(-1),
	}
	return plt
}

// onKey handles keyboard input to modify the plot's value range or toggle the auto-scaling mode.
func (plt *MemPlot) onKey(_ int, key rune) {
	interval := math.Abs(plt.maxVal - plt.minVal)
	scale := (interval * 10) / 100
	switch key {
	case 'a', '+':
		plt.auto = false
		plt.maxVal += scale
		plt.minVal -= scale
	case 'z', '-':
		plt.auto = false
		plt.maxVal -= scale
		plt.minVal += scale
	case 'r':
		plt.auto = !plt.auto
	}
}

// onTimer is a method that periodically collects memory statistics, updates the MemPlot data, and triggers a repaint.
func (plt *MemPlot) onTimer() {
	var rtStats runtime.MemStats
	runtime.ReadMemStats(&rtStats)

	var val float64
	switch plt.kind {
	case 0:
		val = bToMb(rtStats.Alloc)
	case 1:
		val = bToMb(rtStats.TotalAlloc)
	case 2:
		val = bToMb(rtStats.Sys)
	case 3:
		val = float64(rtStats.NumGC)
	default:
		val = bToMb(rtStats.Alloc)
	}
	if val < plt.minVal {
		plt.minVal = val
	}
	if val > plt.maxVal {
		plt.maxVal = val
	}
	plt.container = append(plt.container, val)
	if len(plt.container) > 10 {
		plt.container = plt.container[1:]
	}
	fmt.Println("onTimer plt.container", plt.container)
	//fmt.Println("kernel.PaintRequest()")
}

// onPaint handles the rendering of the plot on the given surface using the current data and value range.
// It adjusts the minimum and maximum plot values if auto-scaling is disabled.
// The method invokes the DrawSeries function of the provided ISurface to display the data.
func (plt *MemPlot) onPaint(surface ISurface) {
	/*
		var minPlot = 0.0
		var maxPlot = 0.0
		if !plt.auto {
			minPlot = plt.minVal
			maxPlot = plt.maxVal
		}
		fmt.Println("onPaint plt.container", i, plt.container, minPlot, maxPlot)
		//surface.DrawSeries(plt.container, -1, -1, minPlot, maxPlot)

	*/

	var minPlot = 0.0
	var maxPlot = 0.0
	if !plt.auto {
		minPlot = plt.minVal
		maxPlot = plt.maxVal
	}

	//var i int = 1000
	//plt.container = append(plt.container, 1.8)
	//if len(plt.container) > 10 {
	//	plt.container = plt.container[1:]
	//}
	surface.DrawSeries(plt.container, -1, -1, minPlot, maxPlot)
	//fmt.Println("onPaint plt.container", surface, plt.container, minPlot, maxPlot)
}

// _instance is a singleton instance of the MemPlot structure, used for memory plotting and rendering operations.
var _instance *MemPlot = NewMemPlot(0)

// onPaint handles the repaint event for the current graphical instance by delegating the operation to the _instance object.
func onPaint(s interface{}) {
	z, _ := s.(ISurface)
	//var z ISurface = s
	//fmt.Println("onPaint", z)
	_instance.onPaint(z)
	//_instance.onPaint(z)
}

// onTimer triggers the onTimer method of the _instance object, typically used for handling periodic actions or events.
func onTimer() {
	//_instance.onKey(0, 'r')
	_instance.onTimer()
}

// main is the entry point of the program that initializes a MemPlot instance based on the first argument and sets up a timer.
func main(args []string) {
	kind := 0
	if len(args) > 0 {
		switch args[0] {
		case "alloc":
			kind = 0
		case "total":
			kind = 1
		case "os":
			kind = 2
		case "gc":
			kind = 3
		}
	}
	//_instance =
	fmt.Println("kernel.CreateTimer(0, 300, -1)", kind)
}
