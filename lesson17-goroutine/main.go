package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func heavyTask(wg *sync.WaitGroup) {
	defer wg.Done()
	sum := 0
	for i := 0; i < 100e8; i++ {
		sum += i
	}
}

func main() {

	var wg sync.WaitGroup

	start := time.Now()

	// numCPU := runtime.NumCPU()
	numCPU := 4

	fmt.Printf("Number of CPU cores: %d\n", numCPU)

	runtime.GOMAXPROCS(numCPU)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			heavyTask(&wg)
		}()
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("All tasks completed in %s\n", elapsed)	

}