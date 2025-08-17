package stub

const Source3 = `
package main

import "fmt"

type Home struct {
	Name string
	Age int
	Address string
}

//func NewHome() Home {
//	return Home{Name: "Mario", Age: 20, Address: "Shanghai"}
//}

func (h *Home) test() {
	h.Name = "Alfa"
	//fmt.Println("home", h)
}

func main() {
	h := Home{Name:"Alfa", Age: 20, Address: "Shanghai"}
	//a:="prova"
	//fmt.Println(a)
	//h := NewHome()
	h.test()
fmt.Println(h)
}
`
