package sources

import "fmt"

func increment(val *int) {
	*val = *val + 1
}

func main() {
	x := 10
	ptrX := &x
	//*ptrX = *ptrX + 1
	increment(ptrX)
	fmt.Println(x)
}
