package stub

const Source1 = `
package main

import "fmt"

type MyStruct struct {
	Name string
}

func (z * MyStruct) PrintA(a string) int {
	fmt.Println(a)
	return 10
}

func (z * MyStruct) PrintB(b string) int {
	fmt.Println(b)
	return 20
}

func testFunc(a int) (string, int) {
	fmt.Println("Hello, world!")
	return "prova", 2
}

func main() {
	//a := test()
	//fmt.Println(a)
	//for x:=10; x==0; x-- {
	//	fmt.Println(x)
	//}
	//fmt.Println("PROVA")
	//var x = 4
	//var y = 15
	//z := x+y
	//b := len("Hello, world!,Hello, world!, Hello,world!,Hello, world!,Hello, world!")
	//fmt.Println(b)
	//fmt.Println(z)
	//return z
	a := 10
    kk := MyStruct{}
	retA := kk.PrintA("prova")
	retB := kk.PrintB("prova")
	myVar, mainRet := testFunc(5)
	fmt.Println(retA, retB, myVar, mainRet)
}
`
