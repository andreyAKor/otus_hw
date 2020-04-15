package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

	wg  *sync.WaitGroup = &sync.WaitGroup{}
	mux *sync.Mutex     = &sync.Mutex{}

	quitCh chan struct{}
	taskCh chan Task

	errs  int
	errCh chan error

	countTasks int
)

type Task func() error

func worker(m int) {
	defer wg.Done()

	for {
		select {
		case <-quitCh:
			return
		default:
		}

		select {
		case <-quitCh:
			return
		case task := <-taskCh:
			if err := task(); err != nil {
				mux.Lock()
				errs++

				if errs == m {
					errCh <- ErrErrorsLimitExceeded
				}

				mux.Unlock()
			}

			mux.Lock()
			countTasks--

			if countTasks == 0 {
				errCh <- nil
			}

			mux.Unlock()
		}
	}
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, n int, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	countTasks = len(tasks)

	errs = 0
	errCh = make(chan error, 1)
	defer close(errCh)

	taskCh = make(chan Task)
	defer close(taskCh)

	defer wg.Wait()

	quitCh = make(chan struct{})
	defer close(quitCh)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker(m)
	}

	for _, task := range tasks {
		select {
		case err := <-errCh:
			return err
		default:
		}

		taskCh <- task
	}

	return <-errCh
}
