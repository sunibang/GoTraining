package numbers

import "fmt"

func SumAndPrint(numbers []int) {
	result := 0
	for _, number := range numbers {
		result += number
	}
	fmt.Println("Result:", result)
}
