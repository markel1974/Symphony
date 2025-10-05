package stub

import (
	"fmt"
	"sort"
	"time"
)

type Plotter struct {
}

func (p *Plotter) DrawSeries(data []float64, rows int, columns int, min float64, max float64) {
	fmt.Println("Plotter draw series", data, rows, columns, min, max)
}

func MemPlotTest(debug bool) error {
	var files []string
	err := walk(".", "mem_plot", &files)
	if err != nil {
		return nil
	}
	if len(files) != 1 {
		return fmt.Errorf("expected 1 file, got %d", len(files))
	}
	fi := NewFileInfo(files[0])
	l := NewLauncher(debug)
	if err = l.Setup("native", []*FileInfo{fi}); err != nil {
		return err
	}
	mainArgs := []interface{}{[]string{"os"}}
	_, err = l.Exec("main", mainArgs)
	if err != nil {
		return err
	}
	onPaintArgs := []interface{}{&Plotter{}}
	var onTimerArgs []interface{} = nil

	for x := 0; x < 15; x++ {
		_, err = l.Exec("onTimer", onTimerArgs)
		if err != nil {
			return err
		}
		_, err = l.Exec("onPaint", onPaintArgs)
		if err != nil {
			return err
		}
	}
	return nil
}

// Launch searches for files with the given prefix, sorts them, and executes the specified entry point with arguments.
// It uses a setup and execution workflow, returning results or an error if processing fails.
func Launch(prefix string, seqId string, entry string, args []interface{}, debug bool) error {
	var files []string
	err := walk(".", prefix, &files)
	if err != nil {
		return nil
	}
	sort.Strings(files)

	//gk := objects.NewGateKeeper()
	for _, fileName := range files {
		fmt.Printf("\n\n------------------ %s ------------------\n", fileName)
		fi := NewFileInfo(fileName)
		l := NewLauncher(debug)
		if err = l.Setup(seqId, []*FileInfo{fi}); err != nil {
			return err
		}
		start := time.Now()
		rv, err := l.Exec(entry, args)
		if err != nil {
			return err
		}
		end := time.Since(start)
		_, allocatedObjects, iterations, maxFrames := l.vm.Statistics()
		fmt.Println("------------- Result -------------")
		fmt.Printf("elapsed: %v, allocatedObjects: %d, iterations: %d, frames: %d\n", end, allocatedObjects, iterations, maxFrames)
		fmt.Printf("return values:%v\n", rv)
	}
	return nil
}
