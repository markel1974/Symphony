package stub

import "fmt"

type Home struct {
	Name    string
	Age     int
	Address string
}

func Beta() {
	fmt.Println("test", 5+6)
}

func NewHome() *Home {
	return &Home{Name: "Mario", Age: 20 + 5, Address: "Shanghai"}
}

func (h *Home) test() {
	h.Name = "Alfio"
	//z := h.Age + 5
	//h.Age = z
	h.Age = h.Age + 5
	fmt.Println("home", h.Name, h.Age, h.Address)
}

func main() {
	//h := Home{Name:"Alfa", Age: 20, Address: "Shanghai"}
	//a:="prova"
	//fmt.Println(a)
	home := NewHome()
	home.test()
	Beta()
}
