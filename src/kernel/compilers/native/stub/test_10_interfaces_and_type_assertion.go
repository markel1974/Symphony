package stub

import "fmt"

type Shape interface {
	Area() int
}

type RectangleShape struct {
	width  int
	height int
}

func (r RectangleShape) Area() int {
	return r.width * r.height
}

type CircleShape struct {
	radius int
}

func (c CircleShape) Area() int {
	return 3 * c.radius * c.radius
}

func getTotalArea(s Shape) int {
	return s.Area()
}

func main() {
	fmt.Println("--- Running Test 10: Interfaces and Type Assertion ---")

	r := RectangleShape{width: 10, height: 4}
	c := CircleShape{radius: 5}

	var s Shape
	s = r
	total := getTotalArea(s)

	s = c
	total = total + getTotalArea(s)

	var shapeVar Shape = RectangleShape{width: 3, height: 3}
	rect, ok := shapeVar.(RectangleShape)

	finalValue := 0
	if ok {
		finalValue = total + rect.Area()
	} else {
		finalValue = -1 // Forza il fallimento se la type assertion non va a buon fine
	}

	expectedValue := 124

	if finalValue == expectedValue {
		fmt.Println("[TEST PASSED] Interfaces, polymorphism, and type assertion worked correctly.")
	} else {
		fmt.Println("[TEST FAILED] Error in interface or type assertion logic.")
	}
}
