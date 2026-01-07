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

	m := make(map[string]int)
	m["TEST"] = 1234

	m1 := map[string]interface{}{
		"TEST": map[string]interface{}{
			"INNER1": map[string]interface{}{
				"INNER2": 30,
			},
		},
	}

	x := &Outer{Name: "Alfa", Value: 100, C: Center{Central: "Beta", I: Inner{Internal: 21, TEST: []float64{1.1, 2.1, 3.1}}}}
	y := &Outer{}
	y.Value += c * 2
	y.C.I.Internal = c * 2
	y.C.Central = "Y"

	x.C.I.Internal = c

	fmt.Println("-------- X")
	fmt.Println(x.C.I.Internal, x.C.Central, x.Value, x.C.I.TEST)
	fmt.Println(x.C.Central)

	fmt.Println("-------- Y")
	fmt.Println(y.C.I.Internal, y.C.Central, y.Value, y.C.I.TEST)
	fmt.Println(y.C.Central)

	fmt.Println("-------- M")
	fmt.Println(m)

	fmt.Println("-------- M1")
	fmt.Println(m1)
	fmt.Println(m1["TEST"])
	fmt.Println(m1["TEST"].(map[string]interface{})["INNER1"].(map[string]interface{})["INNER2"])
}
