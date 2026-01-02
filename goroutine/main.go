package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

var (
	ErrGoroutinePoolStop = errors.New("goroutine pool is stopped")
	ErrUserCancel        = errors.New("user cancelled")
	ErrTimeout           = errors.New("task timeout")
	ErrTaskPanic         = errors.New("task panic")
)

type FutureData struct {
	Err    error
	Result any
}

type Future struct {
	TaskID string
	Result chan FutureData
	cancel context.CancelFunc
}

func (f *Future) Cancel() {
	if f.cancel != nil {
		f.cancel()
	}
}

type GoroutinePoolConfig struct {
	WorkerNum int
	QueueSize int
}

type GoroutinePool struct {
	num           int
	isStop        int32
	closeChan     chan struct{}
	taskChan      chan func()
	wg            sync.WaitGroup
	activeWorkers int32
	totalTasks    int64
}

func NewGoroutinePool(config GoroutinePoolConfig) *GoroutinePool {
	if config.WorkerNum <= 0 {
		config.WorkerNum = 10
	}
	if config.QueueSize <= 0 {
		config.QueueSize = config.WorkerNum * 2
	}

	return &GoroutinePool{
		num:       config.WorkerNum,
		closeChan: make(chan struct{}),
		taskChan:  make(chan func(), config.QueueSize),
	}
}

func (g *GoroutinePool) Enqueue(ctx context.Context,
	task func(context.Context) (any, error),
	timeout time.Duration) (*Future, error) {

	if atomic.LoadInt32(&g.isStop) == 1 {
		return nil, ErrGoroutinePoolStop
	}

	taskCtx, cancel := context.WithCancel(ctx)
	if timeout > 0 {
		taskCtx, cancel = context.WithTimeout(ctx, timeout)
	}

	result := &Future{
		TaskID: uuid.New().String(),
		Result: make(chan FutureData, 1),
		cancel: cancel,
	}

	g.wg.Add(1)
	atomic.AddInt64(&g.totalTasks, 1)

	qTask := func() {
		defer func() {
			g.wg.Done()
			cancel()
			atomic.AddInt32(&g.activeWorkers, -1)

			if r := recover(); r != nil {
				result.Result <- FutureData{
					Err: fmt.Errorf("%w: %v", ErrTaskPanic, r),
				}
			}
		}()

		atomic.AddInt32(&g.activeWorkers, 1)

		done := make(chan FutureData, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					done <- FutureData{
						Err: fmt.Errorf("%w: %v", ErrTaskPanic, r),
					}
				}
			}()
			taskResult, err := task(taskCtx)
			done <- FutureData{Result: taskResult, Err: err}
		}()

		select {
		case data := <-done:
			result.Result <- data
		case <-taskCtx.Done():
			result.Result <- FutureData{Err: ErrTimeout}
		}
	}

	select {
	case g.taskChan <- qTask:
		return result, nil
	case <-g.closeChan:
		g.wg.Done()
		cancel()
		return nil, ErrGoroutinePoolStop
	case <-ctx.Done():
		g.wg.Done()
		cancel()
		return nil, ErrUserCancel
	}
}

func (g *GoroutinePool) Run() {
	for i := 0; i < g.num; i++ {
		go func(workerID int) {
			for {
				select {
				case <-g.closeChan:
					return
				case task := <-g.taskChan:
					task()
				}
			}
		}(i)
	}
}

func (g *GoroutinePool) Stop() {
	if atomic.CompareAndSwapInt32(&g.isStop, 0, 1) {
		close(g.closeChan)
		// 不关闭 taskChan，让已入队的任务执行完
	}
}

func (g *GoroutinePool) Wait() {
	g.wg.Wait()
}

func (g *GoroutinePool) Stats() (active int32, total int64, pending int) {
	return atomic.LoadInt32(&g.activeWorkers),
		atomic.LoadInt64(&g.totalTasks),
		len(g.taskChan)
}

func main() {
	pools := NewGoroutinePool(GoroutinePoolConfig{
		WorkerNum: 10,
		QueueSize: 50,
	})
	pools.Run()
	defer pools.Stop()

	results := make([]*Future, 0, 100)
	for i := 0; i < 100; i++ {
		data := i
		result, err := pools.Enqueue(
			context.Background(),
			func(ctx context.Context) (any, error) {
				result := 0.0
				for j := data; j < 100000; j++ {
					// 检查是否被取消
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					default:
					}
					result += math.Pow(float64(data), 2)
				}
				return result, nil
			},
			5*time.Second,
		)
		if err != nil {
			fmt.Println("enqueue failed, err:", err)
			continue
		}
		results = append(results, result)
	}

	// 等待所有任务完成
	for _, result := range results {
		data := <-result.Result
		if data.Err != nil {
			fmt.Printf("%s: error: %v\n", result.TaskID, data.Err)
		} else {
			fmt.Printf("%s: result: %v\n", result.TaskID, data.Result)
		}
	}

	// 打印统计信息
	active, total, pending := pools.Stats()
	fmt.Printf("\nStats - Active: %d, Total: %d, Pending: %d\n", active, total, pending)

	pools.Wait()
	fmt.Println("All tasks completed")
}
