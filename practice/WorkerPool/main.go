package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	ID int
}

type Result struct {
	JobID    int
	WorkerID int
	Duration time.Duration
}

func Worker(id int, jobChan <-chan Job, resultChan chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobChan {
		start := time.Now()

		// Random thời gian xử lý mỗi job: từ 50ms đến 300ms
		processTime := time.Duration(rand.Intn(250)+50) * time.Millisecond
		time.Sleep(processTime)

		resultChan <- Result{
			JobID:    job.ID,
			WorkerID: id,
			Duration: time.Since(start),
		}
	}
}

func main() {
	const (
		totalJobs  = 100
		numWorkers = 5
	)

	start := time.Now()

	jobChan := make(chan Job, totalJobs)
	resultChan := make(chan Result, totalJobs)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for w := 1; w <= numWorkers; w++ {
		go Worker(w, jobChan, resultChan, &wg)
	}

	go func() {
		for j := 1; j <= totalJobs; j++ {
			jobChan <- Job{ID: j}
		}
		close(jobChan)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Thống kê theo worker: số job xử lý + tổng thời gian
	type stat struct {
		count int
		total time.Duration
	}
	workerStats := make(map[int]*stat)

	for result := range resultChan {
		if workerStats[result.WorkerID] == nil {
			workerStats[result.WorkerID] = &stat{}
		}
		workerStats[result.WorkerID].count++
		workerStats[result.WorkerID].total += result.Duration
	}

	fmt.Println("===== THỐNG KÊ THEO WORKER =====")
	for w := 1; w <= numWorkers; w++ {
		s := workerStats[w]
		avg := s.total / time.Duration(s.count)
		fmt.Printf("Worker%d: %3d job | tổng %v | trung bình %v/job\n", w, s.count, s.total, avg)
	}

	fmt.Printf("\n===== Hoàn thành %d job bằng %d worker sau %v =====\n", totalJobs, numWorkers, time.Since(start))
}
