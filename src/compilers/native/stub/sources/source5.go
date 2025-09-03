package sources

import "fmt"

func test5() int {
	return 10
}

func main(a int, b int) int {
	//var input string = "INPUT"
	//results := "RESULTS"
	//for y:=10; y==0; y-- {
	//	fmt.Printf("y:%d -> %s\n", y, "CIAO")
	//}
	y := len("CIAO")
	k := test5()
	fmt.Printf("y:%d -> %s %d\n", y, "CIAO", k)
	fmt.Printf("y:%d -> %s %d\n", y, "CIAO", test5())
	for y := 0; y < 10; y++ {
		fmt.Printf("y:%d -> %s\n", y, "CIAO")
	}
	//z := a + b
	//y := a - b
	//x := a * b
	//fmt.Println(input, a, b, results, z, z, y, x)
	return 10
}
