package stub

import "fmt"

func Beta4() (int, bool) {
	return 10, false
}

func main() {
	a, b := Beta4()
	fmt.Println("PROVA", a, b)
}
