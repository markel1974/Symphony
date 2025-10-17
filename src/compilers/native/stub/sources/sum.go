package sources

import "fmt"

type Inner struct {
	Internal int
}

type Center struct {
	Central string
	I       Inner
}
type Outer struct {
	Name  string
	Value int64
	C     Center
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
