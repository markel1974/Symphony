// //go:build no_compile
package sources

import (
	"fmt"
	"math"
	"runtime"
	//"kernel"
)

// ISurface represents an abstraction for rendering data-driven graphical series onto a surface.
// DrawSeries renders a series of data points within specified boundaries and dimensions.
type ISurface interface {
	DrawSeries(data []float64, rows int, columns int, min float64, max float64)
}

type Surface struct {
}

func (s *Surface) DrawSeries(data []float64, rows int, columns int, min float64, max float64) {
	fmt.Println("Drawing series:", data, rows, columns, min, max)
}

// bToMb converts a given size in bytes (b) to megabytes as a float64.
func byteToMegaByte(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

// MemPlot is a structure that tracks memory statistics over time for visualization and analysis.
type MemPlot struct {
	data   []float64
	kind   int
	minVal float64
	maxVal float64
	auto   bool
}

// NewMemPlot initializes a new MemPlot instance of the specified kind, with auto scaling enabled and empty data.
func NewMemPlot(kind int) *MemPlot {
	plt := &MemPlot{
		kind:   kind,
		auto:   true,
		minVal: math.Inf(1),
		maxVal: math.Inf(-1),
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
	var ms runtime.MemStats
	z := &ms
	runtime.ReadMemStats(z)
	fmt.Println(z)
	fmt.Println("Plotting:", plt)
	//kernel.PaintRequest()
}

// onPaint handles the rendering of the plot on the given surface using the current data and value range.
// It adjusts the minimum and maximum plot values if auto-scaling is disabled.
// The method invokes the DrawSeries function of the provided ISurface to display the data.
func (plt *MemPlot) onPaint(surface ISurface) {
	var minPlot float64 = 0
	var maxPlot float64 = 0
	if !plt.auto {
		minPlot = plt.minVal
		maxPlot = plt.maxVal
	}
	fmt.Println("Painting:", minPlot, maxPlot)
	surface.DrawSeries(plt.data, -1, -1, minPlot, maxPlot)
}

// _instance is a singleton instance of the MemPlot structure, used for memory plotting and rendering operations.
var _instanceMemPlot *MemPlot //NewMemPlot(0)

// onPaint handles the repaint event for the current graphical instance by delegating the operation to the _instance object.
func onPaintEntry() {
	var i ISurface
	i = &Surface{}
	_instanceMemPlot.onPaint(i)
}

// onTimer triggers the onTimer method of the _instance object, typically used for handling periodic actions or events.
func onTimerEntry() {
	_instanceMemPlot.onTimer()
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
	_instanceMemPlot = NewMemPlot(kind)
	//fmt.Println("KIND", kind)
	//fmt.Println("INSTANCE", _instance)

	//onPaintEntry()

	onTimerEntry()
	//kernel.CreateTimer(0, 300, -1)
}
