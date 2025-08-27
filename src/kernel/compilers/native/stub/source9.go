package stub

import "fmt"

const TEST = "ALFA_BETA"

func TaskFilterGenerator9(z int) func(a int) int {
	return func(int) int {
		return z * 2
	}
}

func main() {

	fmt.Println(TEST)
	filterIncomplete := TaskFilterGenerator9
	allTasks := []int{1, 2, 3}
	for _, taskT := range allTasks {
		z := filterIncomplete(taskT)
		result := z(5)
		fmt.Println(result)
	}
}
