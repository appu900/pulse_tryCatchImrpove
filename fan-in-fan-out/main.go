package main

import (
	"fmt"
	"sync"
	"time"
)

func fanOut(input <-chan int, numWorkers int) []<-chan int {
	outputs := make([]<-chan int, numWorkers)
	for i := 0; i < numWorkers; i++ {
		fmt.Println("Process stared to create a queue with id:",i)
		outputs[i] = worker(i, input)

	}
	return outputs
}

func worker(id int, input <-chan int) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		for num := range input {
			time.Sleep(time.Duration(100*(id+1)) * time.Millisecond)
			result := num * num
			fmt.Printf("worker %d processed %d -> %d\n", id, num, result)
			output <- result
			fmt.Println("Finished processing create a queue with id", id)
		}
	}()
	return output

}

func fanIn(inputs []<-chan int) <-chan int {
	output := make(chan int)
	var wg sync.WaitGroup
	for _, ch := range inputs {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for val := range ch {
				output <- val
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(output)
	}()
	return output
}

func main() {
	inputChannal := make(chan int)
	go func() {
		defer close(inputChannal)
		for i := 0; i < 10; i++ {
			inputChannal <- i
		}
	}()

	workers := fanOut(inputChannal, 3)
	results := fanIn(workers)
	var resultSlice []int
	for result := range results {
		resultSlice = append(resultSlice, result)
	}

	fmt.Println("All results:", resultSlice)
}
