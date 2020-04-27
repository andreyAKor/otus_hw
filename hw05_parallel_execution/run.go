package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, n int, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	wg := &sync.WaitGroup{}
	errs := int64(0)
	taskCh := make(chan Task)
	doneCh := make(chan struct{})

	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker(m, wg, doneCh, taskCh, &errs)
	}

	for _, task := range tasks {
		select {
		case <-doneCh:
			wg.Wait()
			return ErrErrorsLimitExceeded
		case taskCh <- task:
		}
	}

	close(taskCh)
	wg.Wait()

	select {
	case <-doneCh:
		return ErrErrorsLimitExceeded
	default:
	}

	return nil
}

func worker(
	m int,
	wg *sync.WaitGroup,
	doneCh chan struct{},
	taskCh chan Task,
	errs *int64,
) {
	defer wg.Done()

	for {
		select {
		case <-doneCh:
			return
		default:
		}

		select {
		case <-doneCh:
			return
		case task, ok := <-taskCh:
			if !ok {
				return
			}

			if err := task(); err != nil {
				if atomic.AddInt64(errs, 1) == int64(m) {
					close(doneCh)
				}
			}
		}
	}
}
