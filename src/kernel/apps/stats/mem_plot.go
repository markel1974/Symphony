package stats

import (
	"math"
	"runtime"

	"github.com/markel1974/symphony/src/kernel/interfaces"
	"github.com/markel1974/symphony/src/kernel/process"
)

// bToMb converts a byte value (uint64) to megabytes (float64) by dividing the input by 1024 twice.
func bToMb(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

// MemPlot represents a memory plotting utility for capturing and visualizing runtime memory statistics in real-time.
type MemPlot struct {
	process interfaces.IUserProcess
	data    []float64
	kind    int
	minVal  float64
	maxVal  float64
	auto    bool
}

// NewMemPlot creates a new instance of MemPlot with the specified kind, initializing its fields with default values.
func NewMemPlot(kind int) *MemPlot {
	plt := &MemPlot{
		kind:   kind,
		auto:   true,
		data:   nil,
		minVal: math.Inf(1),
		maxVal: math.Inf(-1),
	}
	return plt
}

// Setup initializes the MemPlot instance by setting up the required process handlers for key, timer, and paint events.
func (plt *MemPlot) Setup(process interfaces.IUserProcess) {
	plt.process = process
	plt.process.SetOnKey(plt.onKey)
	plt.process.SetOnTimer(plt.onTimer)
	plt.process.SetOnPaint(plt.onPaint)
}

// Start initializes the timer for the MemPlot instance with a 300 millisecond interval and an infinite execution count.
func (plt *MemPlot) Start() {
	plt.process.CreateTimer(0, 300, -1)
}

// onKey handles keyboard input events, allowing adjustments to plot scaling or toggling the auto-scaling feature.
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

// onTimer is invoked periodically based on the timer to collect and track memory statistics and update the plot data.
func (plt *MemPlot) onTimer(_ int, _ int) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
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
	plt.process.PaintRequest()
}

// onPaint renders the graph data on the provided surface, using custom or automatic plot range values.
func (plt *MemPlot) onPaint(surface interfaces.ISurface) {
	var minPlot float64 = 0
	var maxPlot float64 = 0
	if !plt.auto {
		minPlot = plt.minVal
		maxPlot = plt.maxVal
	}
	surface.DrawSeries(plt.data, -1, -1, minPlot, maxPlot)
}

// CreateMemPlot creates a new runtime memory plotting command, initializing its setup and starting the memory monitoring process.
func CreateMemPlot() interfaces.ICommand {
	run := func(process interfaces.IUserProcess, args []string) error {
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
		plt := NewMemPlot(kind)
		plt.Setup(process)
		plt.Start()
		return nil
	}
	root := process.NewCommand("rtplot", interfaces.CommandTypeFile, nil, true, run)
	root.SetHelp("Runtime Plot", "Runtime Plot")
	return root
}
