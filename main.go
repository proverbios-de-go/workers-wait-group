package main

import (
	"fmt"
	"sync"
	"time"
)

type jobItem struct {
	jobID int
}

type resultItem struct {
	jobID int
}

func worker(wg *sync.WaitGroup, jobs chan jobItem, results chan resultItem) {
	defer wg.Done()
	for job := range jobs {
		fmt.Println(fmt.Sprintf("Working on new job:%d", job.jobID))
		time.Sleep(3 * time.Second)
		results <- resultItem{
			jobID: job.jobID,
		}
	}

}

func main() {

	jobsC := make(chan jobItem)
	resultsC := make(chan resultItem)
	wg := new(sync.WaitGroup)

	for i := 1; i <= 2; i++ {
		//fmt.Println(fmt.Sprintf("Bootstraping worker:%d", i))
		wg.Add(1)
		go worker(wg, jobsC, resultsC)
	}

	var jobs []jobItem
	for i := 1; i <= 10; i++ {
		jobs = append(jobs, jobItem{jobID: i})
		fmt.Println(fmt.Sprintf("Creating new job, all jobs:%v", jobs))
	}

	go func() {
		for _, job := range jobs {
			fmt.Println(fmt.Sprintf("Sending job:%d", job.jobID))
			jobsC <- job
		}
		close(jobsC)
	}()

	go func() {
		wg.Wait()
		close(resultsC)
	}()

	for result := range resultsC {
		fmt.Println(fmt.Sprintf("This is the result of the jobID:%d", result.jobID))
	}

	fmt.Println("End")
}
