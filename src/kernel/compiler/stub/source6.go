package stub

const Source6 = `
package main

import "fmt"

type TestResult struct {
    Name    string
    Passed  bool
    Message string
}

//var testResults []TestResult = []TestResult{ { "test1", true, "ok" }, { "test2", false, "error" }}
var testResults = TestResult{ "test2", false, "error" }

func main(a int, b int) {
	fmt.Println(testResults)
}
`
