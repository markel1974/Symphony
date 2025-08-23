package stub

const Source6 = `
package main

import "fmt"

type TestResult struct {
    Name    string
    Passed  bool
}

//var testResults []TestResult = []TestResult{ { "var1", true}, { "var2", false }}
var testResults = TestResult{ "var1", false }
//var testResults = 1
func main() {
	//var testResults = TestResult{ "var1", false }
	testResults.Name = "prova"

	fmt.Println(testResults)
}
`
