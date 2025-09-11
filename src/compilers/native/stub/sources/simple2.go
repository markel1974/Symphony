package sources

import (
	"fmt"
	"runtime"
)

func main() {
	varMemStats := runtime.MemStats{}
	runtime.ReadMemStats(&varMemStats)
	defer fmt.Println("Hello", varMemStats.BuckHashSys)
	//defer fmt.Println("Hello2")
	//fmt.Println("BuckHashSys", varMemStats.BuckHashSys)
	fmt.Println("Alloc", varMemStats.Alloc)
}
