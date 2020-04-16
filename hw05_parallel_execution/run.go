package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

// ErrErrorsLimitExceeded result Run
var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

// Task interface for Run
type Task func() error

// Running N gorutin. The function returns a channel signaling the completion of all goroutines
func runGorutins(count int, taskCh <-chan Task, errorCh chan<- error) <-chan struct{} {
	result := make(chan struct{})
	go func() {
		var waiting sync.WaitGroup
		waiting.Add(count)
		for i := 0; i < count; i++ {
			// Run gorutin
			go func() {
				defer waiting.Done()
				for t := range taskCh {
					errorCh <- t()
				}
			}()
		}
		waiting.Wait()
		// All goroutines are completed
		close(result)
	}()
	return result
}

var errorCounter = func(maxCountErrors int) func(error) bool {
	errorCount := 0
	return func(e error) bool {
		if maxCountErrors < 0 {
			return true
		}
		if errorCount >= maxCountErrors {
			return false
		}
		if e != nil {
			errorCount++
		}
		if errorCount >= maxCountErrors {
			return false
		}
		return true
	}
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, n int, m int) (errorResult error) {
	errorCh := make(chan error)
	taskCh := make(chan Task)

	safeCloseTasks := func() func() {
		isClosed := false
		return func() {
			if isClosed {
				return
			}
			close(taskCh)
			isClosed = true
		}
	}()
	waitingCh := runGorutins(n, taskCh, errorCh)
	checkError := errorCounter(m)

	indexTask := 0
	countTasks := len(tasks)
	// error handling. Close channel when error limit is reached
	handlerError := func(e error) {
		if !checkError(e) {
			if errorResult == nil {
				errorResult = ErrErrorsLimitExceeded
				indexTask = countTasks
				safeCloseTasks()
			}
		}
	}
	for {
		if indexTask < countTasks {
			// Send tasks and process errors
			select {
			case taskCh <- tasks[indexTask]:
				indexTask++
				if indexTask >= countTasks {
					safeCloseTasks()
				}
			case e := <-errorCh:
				handlerError(e)
			case <-waitingCh:
				return errorResult
			}
			continue
		}
		// Process errors from goroutines
		select {
		case e := <-errorCh:
			handlerError(e)
		case <-waitingCh:
			return errorResult
		}
	}
}
