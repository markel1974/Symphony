package sources

import "fmt"

type Root struct {
	Name  string
	Value int64
}

func main() {
	//a := 20
	//a /= 10
	//// a = a + 10
	//fmt.Println(a)

	const c = 31
	r := &Root{Name: "Alfa", Value: 100}
	//r.Value += c
	r.Value = r.Value + c
	fmt.Println(r.Value)
}
