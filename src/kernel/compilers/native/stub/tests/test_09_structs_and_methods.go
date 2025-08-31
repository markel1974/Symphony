package tests

import "fmt"

type RectangleStruct struct {
	width  int
	height int
}

func (r RectangleStruct) Area() int {
	return r.width * r.height
}

func NewRectangleStruct(width int, height int) RectangleStruct {
	return RectangleStruct{width: width, height: height}
}

func main() {
	fmt.Println("--- Running Test: Structs and Methods ---")

	r1 := RectangleStruct{width: 10, height: 5}
	area := r1.Area()
	sumFields := r1.width + r1.height

	finalValue := area + sumFields
	expectedValue := 65

	if finalValue == expectedValue {
		fmt.Println("[TEST PASSED] Structs and methods worked correctly.")
	} else {
		fmt.Println("[TEST FAILED] Error in struct or method handling.")
	}
}
