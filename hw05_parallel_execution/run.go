package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

// ErrErrorsLimitExceeded result Run
var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

// Task type of argument for Run
type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, n int, m int) (errorResult error) {
	errorCh := make(chan error)
	taskCh := make(chan Task)
	doneCh := make(chan struct{})
	var wg sync.WaitGroup
	// Workers start
	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker(taskCh, errorCh, doneCh, &wg)
	}
	// Run generator of tasks
	wg.Add(1) //nolint:gomnd
	go func() {
		defer wg.Done()
		defer close(taskCh)
		for _, t := range tasks {
			select {
			case taskCh <- t:
			case <-doneCh:
				return
			}
		}
	}()
	// Run handler of errors
	doneHandlerErrors := func() <-chan struct{} {
		ch := make(chan struct{})
		go func() {
			defer close(ch)
			countError := 0
			totalCount := 0
			maxCount := m + n
			for e := range errorCh {
				if m < 0 {
					continue
				}
				totalCount++
				if e != nil {
					countError++
				} else {
					continue
				}
				if (totalCount + countError) >= maxCount {
					errorResult = ErrErrorsLimitExceeded
					close(doneCh)
					return
				}
			}
		}()
		return ch
	}()

	wg.Wait()
	close(errorCh)
	<-doneHandlerErrors
	return
}

// worker handler of tasks
func worker(taskCh <-chan Task, errorCh chan<- error, doneCh <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case t, ok := <-taskCh:
			if !ok {
				return
			}
			e := t()
			select {
			case errorCh <- e:
			case <-doneCh:
				return
			}
		case <-doneCh:
			return
		}
	}
}
