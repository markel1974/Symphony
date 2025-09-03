package sources

import "fmt"

const TEST = "ALFA_BETA"

func TaskFilterGenerator9(z int) func(a int) int {
	return func(int) int {
		return z * 2
	}
}

func main() {
	for y := 0; y < 10; y++ {
		fmt.Print(y)
		fmt.Print(" -> ")
		switch y {
		case 1:
			fmt.Println("i'm 1")
		case 2:
			fmt.Println("i'm two")
		case 3:
			fmt.Println("i'm 3")
		case 4:
			fmt.Println("i'm four")
		default:
			fmt.Println("i'm default")
		}
	}
}
