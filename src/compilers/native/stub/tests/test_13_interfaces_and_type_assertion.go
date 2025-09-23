package tests

import "fmt"

type IShape interface {
	Area() int
}

type RectangleShape struct {
	width  int
	height int
}

func (rs RectangleShape) Area() int {
	return rs.width * rs.height
}

type CircleShape struct {
	radius int
}

func (cs CircleShape) Area() int {
	return 3 * cs.radius * cs.radius
}

func getTotalArea(s IShape) int {
	return s.Area()
}

func main() {
	fmt.Println("--- Running Test: Interfaces and Type Assertion ---")

	rShape := RectangleShape{width: 10, height: 4}
	c := CircleShape{radius: 5}

	var s IShape
	s = rShape
	total := getTotalArea(s)

	s = c
	total = total + getTotalArea(s)

	var shapeVar IShape = RectangleShape{width: 3, height: 3}
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
		fmt.Printf("[TEST FAILED] Error in interface or type assertion logic: %d | %d\n", finalValue, expectedValue)
	}
}
