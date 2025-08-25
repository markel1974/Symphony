package stub

const Source8 = `
package main

import "fmt"

type Task struct {
Name    string
}

func main() {
	t1 := Task{Name: "Task1"}
	t2 := Task{Name: "Task2"}
	t1.Name = "Task1_mod"
	allTasks := []Task{t1, t2}
	for _, taskT := range allTasks {
		fmt.Println(taskT)
	}
}
`
