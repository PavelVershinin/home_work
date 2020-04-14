package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrNCanNotBeLessThanOne = errors.New("N can't be less than one") //nolint:stylecheck

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
//nolint:gocritic
func Run(tasks []Task, N int, M int) error {
	if N < 1 {
		return ErrNCanNotBeLessThanOne
	}

	var wait sync.WaitGroup
	var errCount int32
	limitCh := make(chan struct{}, N)

	for _, task := range tasks {
		limitCh <- struct{}{}

		if int(atomic.LoadInt32(&errCount)) >= M {
			break
		}

		wait.Add(1)
		go func(t Task) {
			defer func() {
				<-limitCh
				wait.Done()
			}()
			if err := t(); err != nil {
				atomic.AddInt32(&errCount, 1)
			}
		}(task)
	}

	wait.Wait()

	if int(errCount) >= M {
		return ErrErrorsLimitExceeded
	}

	return nil
}
