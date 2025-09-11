package sources

import (
	"fmt"
	"runtime"
)

func main() {
	varMemStats0 := runtime.MemStats{}
	runtime.ReadMemStats(&varMemStats0)
	varMemStats1 := runtime.MemStats{}
	runtime.ReadMemStats(&varMemStats1)
	v0 := "0"
	v1 := "1"
	V2 := "2"
	defer fmt.Println(v0, varMemStats0.Sys, v1, V2)
	//defer fmt.Println("Hello2")
	//fmt.Println("BuckHashSys", varMemStats.BuckHashSys)
	//fmt.Println("Alloc", varMemStats.Alloc)
}
