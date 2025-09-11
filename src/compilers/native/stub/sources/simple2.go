package sources

import (
	"fmt"
	"runtime"
)

func main() {
	varMemStats := runtime.MemStats{}
	runtime.ReadMemStats(&varMemStats)
	guilty := varMemStats.BuckHashSys
	fmt.Println("guilty", guilty)
}
