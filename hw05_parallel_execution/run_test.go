package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return fmt.Errorf("error from task %d", i)
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		result := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, result)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		result := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.Nil(t, result)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("performing tasks one at a time", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, tasksCount)
		result := make([]int, 0, tasksCount)
		expected := make([]int, tasksCount)

		for i := 0; i < tasksCount; i++ {
			n := i
			expected[i] = n
			tasks[i] = func() error {
				result = append(result, n)
				return nil
			}
		}

		require.Equal(t, nil, Run(tasks, 1, 1))
		require.Equal(t, expected, result)
	})

	t.Run("N can't be less than one", func(t *testing.T) {
		tasksCount := 10
		runTasksCount := int32(0)
		tasks := make([]Task, tasksCount)

		for i := 0; i < tasksCount; i++ {
			tasks[i] = func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			}
		}
		require.Equal(t, ErrNCanNotBeLessThanOne, Run(tasks, 0, 1))
		require.Equal(t, int32(0), runTasksCount)
	})

	t.Run("what do we do if M less than or equal to zero? 😕", func(t *testing.T) {
		tasksCount := 10
		runTasksCount := int32(0)
		tasks := make([]Task, tasksCount)

		for i := 0; i < tasksCount; i++ {
			tasks[i] = func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			}
		}

		require.Equal(t, ErrErrorsLimitExceeded, Run(tasks, 1, 0))
		require.Equal(t, int32(0), runTasksCount)
	})
}
