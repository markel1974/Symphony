package sources

import "fmt"

type Inner struct {
	Internal int
	TEST     []float64
}

type Center struct {
	I       Inner
	Central string
}
type Outer struct {
	C     Center
	Name  string
	Value int64
}

func main() {
	//a := 20
	//a /= 10
	//// a = a + 10
	//fmt.Println(a)

	const c = 31
	r := &Outer{Name: "Alfa", Value: 100, C: Center{Central: "Beta", I: Inner{Internal: 21}}}
	//r.Value += c
	//r.Value = r.Value + c
	r.C.I.Internal = c
	fmt.Println(r.C.I.Internal)
}
