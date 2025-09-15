package sources

import (
	"fmt"
	"math"
	"runtime"
)

// bToMb converts bytes to megabytes by dividing the input value in bytes by 1024 twice.
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
		data:   []float64{},
		minVal: math.Inf(1),
		maxVal: math.Inf(-1),
	}
	return plt
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

type Memory struct {
	Alloc      uint64
	TotalAlloc uint64
	Sys        uint64
	NumGC      uint32
}

// timerFn is a timer callback method that reads memory stats, updates min/max values, appends data, and triggers repaint.
func (plt *RtPlot) onTimer(_ int, _ int) {
	ms := runtime.MemStats{}
	//out := runtime.ReadMemStats(&ms)
	runtime.ReadMemStats(&ms)
	//m := Memory{}
	//m := runtime.MemStats{}
	//runtime.ReadMemStats(&m)
	var val float64
	switch plt.kind {
	case 0:
		val = bToMb(ms.Alloc)
	case 1:
		val = bToMb(ms.TotalAlloc)
	case 2:
		val = bToMb(ms.Sys)
	case 3:
		val = float64(ms.NumGC)
	default:
		val = bToMb(ms.Alloc)
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
	fmt.Println("Plotting:", plt.data, plt.minVal, plt.maxVal)
	//plt.process.PaintRequest()
}

// paintFn renders the current data series onto the provided ISurface, using defined min and max values or auto-scaling.
func (plt *RtPlot) onPaint() {
	var minPlot float64 = 0
	var maxPlot float64 = 0
	if !plt.auto {
		minPlot = plt.minVal
		maxPlot = plt.maxVal
	}
	fmt.Println("Painting:", minPlot, maxPlot)
}

var _instance = RtPlot{}

func onPaint() (int, int, int, map[string]interface{}) {
	_instance.onPaint()
	z := map[string]interface{}{"1000": "valid", "2000": "valid", "3000": "valid"}
	return 1, 2, 3, z
}

func onTimer(a int, b int) (int, int, int, map[string]interface{}) {
	_instance.onTimer(a, b)
	z := map[string]interface{}{"1000": "valid", "2000": "valid", "3000": "valid"}
	return 3, 2, 1, z
}

func main() (int, int, int, map[string]interface{}) {
	//plt := RtPlot{} //NewRtPlotData(0)
	//_instance.onPaint()
	//test3 := &RtPlot{}
	//test3.minVal += 1
	//_instance.onTimer(0, 0)

	onTimer(0, 0)

	return onPaint()
	//z := map[string]interface{}{"1000": "valid", "2000": "valid", "3000": "valid"}
	//return 3, 2, 1, z
}
