package stub

import "fmt"

type MyStruct struct {
	Name string
}

func (z *MyStruct) PrintA(a string) int {
	fmt.Println(a)
	return 10
}

func (z *MyStruct) PrintB(b string) int {
	fmt.Println(b)
	return 20
}

func testFunc1(a int) (string, int) {
	fmt.Println("Hello, world!")
	return "prova", 2
}

func test1() int {
	return 10
}

func main() {
	b := 4
	var y = 15
	z := b + y

	if z == 19 || b == 4 {
		fmt.Println("PROVA")
	}
	fmt.Println(z)
	//var kk = 5
	//a := test1()
	//fmt.Println(a)
	//for x:=10; x==0; x-- {
	//	fmt.Println(x)
	//}
	//fmt.Println("PROVA")

	//b := len("Hello, world!,Hello, world!, Hello,world!,Hello, world!,Hello, world!")
	//fmt.Println(b)
	//fmt.Println(z)
	//return z
	//a := 10
	//kk := MyStruct{}
	//retA := kk.PrintA("prova")
	//retB := kk.PrintB("prova")
	//myVar, mainRet := testFunc1(5)
	//fmt.Println(retA, retB, myVar, mainRet)
}
