package scheduler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	scheduler := NewSchedulers(10, 10)
	scheduler.Run()

	if !scheduler.running {
		t.Fatal("should running")
	}

	scheduler.Submit(&Task{
		TaskID: "task_1",
		Fn: func(ctx context.Context) (result any, err error) {
			return 123, nil
		},
	})

	result := <-scheduler.Results()
	if result.TaskID != "task_1" {
		t.Fatal("task is should be task_1")
	}

	d, ok := result.Result.(int)
	if !ok {
		t.Fatal("shoule be int")
	}

	if d != 123 {
		t.Fatal("should be 123")
	}

	if scheduler.GetStat().TotalSubmitTasks != 1 {
		t.Fatal("should be 1")
	}

	if scheduler.GetStat().CompleteTasks != 1 {
		t.Fatal("should be 1")
	}
}

func TestScheduler_error(t *testing.T) {
	scheduler := NewSchedulers(10, 10)
	scheduler.Run()

	if !scheduler.running {
		t.Fatal("should running")
	}

	scheduler.Submit(&Task{
		TaskID: "task_1",
		Fn: func(ctx context.Context) (result any, err error) {
			return nil, fmt.Errorf("failed")
		},
	})

	result := <-scheduler.Results()
	if result.TaskID != "task_1" {
		t.Fatal("task is should be task_1")
	}

	if result.Err == nil {
		t.Fatal("shoule be error")
	}

	if scheduler.GetStat().TotalSubmitTasks != 1 {
		t.Fatal("should be 1")
	}

	if scheduler.GetStat().FailedTasks != 1 {
		t.Fatal("should be 1")
	}
}

func TestScheduler_timeout(t *testing.T) {
	scheduler := NewSchedulers(10, 10)
	scheduler.Run()

	if !scheduler.running {
		t.Fatal("should running")
	}

	scheduler.Submit(&Task{
		TaskID:  "task_1",
		Timeout: time.Microsecond,
		Fn: func(ctx context.Context) (result any, err error) {
			resultChan := make(chan TaskResult)
			go func() {
				time.Sleep(time.Second * 10)
				resultChan <- TaskResult{
					Result: nil,
					Err:    fmt.Errorf("failed"),
				}
			}()

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case r := <-resultChan:
				return r.Result, r.Err
			}
		},
	})

	result := <-scheduler.Results()
	if result.TaskID != "task_1" {
		t.Fatal("task is should be task_1")
	}

	if result.Err == nil {
		t.Fatal("shoule be error")
	}

	if scheduler.GetStat().TotalSubmitTasks != 1 {
		t.Fatal("should be 1")
	}

	if scheduler.GetStat().FailedTasks != 1 {
		t.Fatal("should be 1")
	}
}

func TestScheduler_contextcancel(t *testing.T) {
	scheduler := NewSchedulers(2, 10)
	scheduler.Run()

	if !scheduler.running {
		t.Fatal("should running")
	}

	ctx, cancel := context.WithCancel(context.Background())

	task := &Task{
		TaskID: "task_1",
		Fn: func(ctx context.Context) (result any, err error) {
			timer := time.After(time.Second)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-timer:
				return 123, nil
			}
		},
	}

	go func() {
		time.Sleep(time.Microsecond)
		cancel()
	}()

	err := scheduler.SubmitWithContext(ctx, task)
	if err != nil && !errors.Is(err, context.Canceled) {
		t.Fatal("should be cancel")
	}

	result := <-scheduler.Results()
	if !errors.Is(result.Err, context.Canceled) {
		t.Fatalf("should be cancel, cur:%v", result.Err)
	}
}

func TestConcurrencyScheduler(t *testing.T) {
	scheduler := NewSchedulers(10, 10)
	scheduler.Run()

	if !scheduler.running {
		t.Fatal("should running")
	}

	wg := sync.WaitGroup{}
	wg.Go(func() {
		for i := range 10 {
			scheduler.Submit(&Task{
				TaskID: fmt.Sprintf("task_%d", i),
				Fn: func(ctx context.Context) (result any, err error) {
					return i, nil
				},
			})
		}
	})

	wg.Go(func() {
		for i := 10; i < 20; i++ {
			scheduler.Submit(&Task{
				TaskID: fmt.Sprintf("task_%d", i),
				Fn: func(ctx context.Context) (result any, err error) {
					return i, nil
				},
			})
		}
	})

	go func() {
		wg.Wait()
		scheduler.Stop()
	}()

	for result := range scheduler.Results() {
		ids := strings.Split(result.TaskID, "_")
		if len(ids) != 2 {
			t.Fatalf("should be 2")
		}

		if ids[1] != fmt.Sprintf("%d", result.Result) {
			t.Fatalf("should be same, taks_id:%v, result:%v", result.TaskID, result.Result)
		}
	}

	if scheduler.GetStat().TotalSubmitTasks != 20 {
		t.Fatal("should be 20")
	}

	if scheduler.GetStat().CompleteTasks != 20 {
		t.Fatal("should be 20")
	}
}
