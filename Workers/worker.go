package main

import (
	"fmt"
	"sync"
)

type Result struct {
	WorkerID int
	Input    int
	Output   int
}

// The worker function

func worker(id int, data int, process func(int) int, wg *sync.WaitGroup, ch chan Result) {
	defer wg.Done()
	result := process(data)
	ch <- Result{WorkerID: id, Input: data, Output: result}

}

func main() {
	var wg sync.WaitGroup

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// function that processes the data

	process := func(n int) int {
		return n * n
	}

	results := make(chan Result, len(data))

	for i, val := range data {
		wg.Add(1)
		go worker(i, val, process, &wg, results)
	}

	wg.Wait()

	close(results)

	for res := range results {
		fmt.Printf("Go Worker %d processed data %d -> result %d\n", res.WorkerID, res.Input, res.Output)
	}

	fmt.Println("All workers have finished processing via channel.")
}
