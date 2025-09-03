package sources

import (
	"fmt"
	"math"
)

func bToMb(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

// RtPlot represents a real-time data plotting mechanism, managing data, its limits, rendering setup, and auto-scaling.
type RtPlot struct {
	data   []float64
	kind   int
	minVal float64
	maxVal float64
	auto   bool
}

// NewRtPlotData initializes a new RtPlot instance with the specified kind and default settings.
// It sets up the plot to automatically adjust its range and prepares it for plotting runtime data.
// Returns a pointer to the newly created RtPlot instance.
func NewRtPlotData(kind int) *RtPlot {
	plt := &RtPlot{
		kind:   kind,
		auto:   true,
		data:   nil,
		minVal: math.Inf(1),
		maxVal: math.Inf(-1),
	}
	return plt
}

/*
// Setup configures the RtPlot instance with the provided IUserProcess and binds event handlers for read, timer, and paint.
func (plt *RtPlot) Setup(process interfaces.IUserProcess) {
	plt.process = process
	plt.process.SetOnKey(plt.onKey)
	plt.process.SetOnTimer(plt.onTimer)
	plt.process.SetOnPaint(plt.onPaint)
}
*/

// Start initializes and starts a timer for the plot with a zero delay, a 300ms interval, and infinite occurrences.
func (plt *RtPlot) Start() {
	//plt.process.CreateTimer(0, 300, -1)
}

// onKey handles keyboard input, adjusts plot range dynamically, and toggles auto-scaling based on provided key actions.
func (plt *RtPlot) onKey(_ int, key rune) {
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

type MemoryStats struct {
	Alloc      uint64
	TotalAlloc uint64
	Sys        uint64
	NumGC      uint32
}

// timerFn is a timer callback method that reads memory stats, updates min/max values, appends data, and triggers repaint.
func (plt *RtPlot) onTimer(_ int, _ int) {
	var m MemoryStats
	//var m runtime.MemStats
	//runtime.ReadMemStats(&m)
	var val float64
	switch plt.kind {
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
	if val < plt.minVal {
		plt.minVal = val
	}
	if val > plt.maxVal {
		plt.maxVal = val
	}
	plt.data = append(plt.data, val)
	if len(plt.data) > 10 {
		plt.data = plt.data[1:]
	}
	//plt.process.PaintRequest()
}

func (plt *RtPlot) onPaint( /*surface interfaces.ISurface*/ ) {
	var minPlot float64 = 0
	var maxPlot float64 = 0
	if !plt.auto {
		minPlot = plt.minVal
		maxPlot = plt.maxVal
	}
	fmt.Println("onPaint", minPlot, maxPlot)
	//fmt.Println("onPaint", 1, 2)
	//surface.DrawSeries(plt.data, -1, -1, minPlot, maxPlot)
}

func main() {
	plt := NewRtPlotData(0)
	plt.Start()
	plt.onPaint()
}

// paintFn renders the current data series onto the provided ISurface, using defined min and max values or auto-scaling.
//func (plt *RtPlot) onPaint(surface interfaces.ISurface) {
//	var minPlot float64 = 0
//	var maxPlot float64 = 0
//	if !plt.auto {
//		minPlot = plt.minVal
//		maxPlot = plt.maxVal
//	}
//	surface.DrawSeries(plt.data, -1, -1, minPlot, maxPlot)
//}
