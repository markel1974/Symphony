package sources

import "fmt"

type Test struct {
	Name string
}

var global = 20

var ptr = &global

var ptr2 = &Test{Name: "Test"}

func test2() int {
	return 30
}
func main() {
	//a := global
	b := test2()

	fmt.Println(global, b, *ptr, ptr2.Name) //, //, b)
}
