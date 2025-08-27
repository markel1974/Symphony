package stub

import "fmt"

func TaskFilterGenerator9(z int) func(a int) int {
	return func(int) int {
		return z * 2
	}
}

func main() {
	filterIncomplete := TaskFilterGenerator9
	allTasks := []int{1, 2, 3}
	for _, taskT := range allTasks {
		z := filterIncomplete(taskT)
		result := z(5)
		fmt.Println(result)
	}
}
