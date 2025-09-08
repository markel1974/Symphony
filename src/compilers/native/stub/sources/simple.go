package sources

import "fmt"

var global = 20

func test2() int {
	return 30
}
func main() {
	//a := global
	b := test2()

	fmt.Println(b) //, //, b)
}
