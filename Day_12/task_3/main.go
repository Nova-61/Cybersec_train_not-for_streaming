package main 

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func increment(ptr *int) {

	defer wg.Done()
	*ptr++

}

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for i := range numbers {
		wg.Add(1)
		go increment(&numbers[i])
	}
	
	wg.Wait()
	fmt.Println("Numbers:", numbers)
}