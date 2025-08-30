package stub

import "fmt"

type Task struct {
	Name string
}

func main() {
	t2 := Task{Name: "Task2"}
	//t2.Name = "Task1_mod"
	allTasks := []Task{Task{Name: "Task1"}, t2}
	//allTasks := []Task{}
	for _, taskT := range allTasks {
		fmt.Println(taskT)
	}
}
