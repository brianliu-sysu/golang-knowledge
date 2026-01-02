package scheduler

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Task struct {
	TaskID  string
	Timeout time.Duration
	Fn      func(context.Context) (any, error)
	Ctx     context.Context
}

type TaskResult struct {
	TaskID    string
	Err       error
	Result    any
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}

type Scheduler struct {
	// global context for scheduler
	ctx context.Context
	// cancel
	cancel context.CancelFunc
	// workers
	MaxWorkers int
	// work queue
	TaskChan chan *Task
	// result queue
	ResultChan chan *TaskResult
	// waitgroup
	wg sync.WaitGroup
	// mutex
	mu sync.Mutex
	// running
	running bool
	// stats
	stats *SchedulerStat
}

type SchedulerStat struct {
	mu            sync.Mutex
	CompleteTasks int
	FailedTasks   int
	CancelTasks   int

	TotalSubmitTasks int
	ActiveTasks      int
	QueueTasks       int
}

func NewSchedulers(maxWorks, queueSize int) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		ctx:        ctx,
		cancel:     cancel,
		MaxWorkers: maxWorks,
		TaskChan:   make(chan *Task, queueSize),
		ResultChan: make(chan *TaskResult, queueSize),
		stats:      &SchedulerStat{},
	}
}

func (s *Scheduler) Run() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("scheduler is running")
	}

	for i := range s.MaxWorkers {
		s.wg.Add(1)
		go s.work(i)
	}

	s.running = true
	return nil
}

func (s *Scheduler) Submit(task *Task) error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return fmt.Errorf("scheduler is stoped")
	}
	s.mu.Unlock()

	s.beforeSubmitTaks()

	select {
	case s.TaskChan <- task:
		return nil
	case <-s.ctx.Done():
		return fmt.Errorf("scheduler is stoped")
	}
}

func (s *Scheduler) SubmitWithContext(ctx context.Context, task *Task) error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return fmt.Errorf("scheduler is stoped")
	}
	s.mu.Unlock()

	s.beforeSubmitTaks()

	if task.Ctx == nil {
		task.Ctx = ctx
	} else {
		var cancel context.CancelFunc
		task.Ctx, cancel = s.mergeCtx(task.Ctx, ctx)
		defer cancel()
	}

	select {
	case s.TaskChan <- task:
		return nil
	case <-s.ctx.Done():
		return fmt.Errorf("scheduler is stoped")
	case <-ctx.Done():
		s.afterSubmitTaks()
		return ctx.Err()
	}
}

func (s *Scheduler) mergeCtx(ctx1, ctx2 context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-ctx1.Done():
			cancel()
		case <-ctx2.Done():
			cancel()
		}
	}()

	return ctx, cancel
}

func (s *Scheduler) Results() <-chan *TaskResult {
	return s.ResultChan
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	close(s.TaskChan)
	s.wg.Wait()

	close(s.ResultChan)
}

func (s *Scheduler) Shutdown() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	s.cancel()
	s.wg.Wait()

	close(s.TaskChan)
	close(s.ResultChan)
}

func (s *Scheduler) GetStat() *SchedulerStat {
	s.stats.mu.Lock()
	defer s.stats.mu.Unlock()

	return &SchedulerStat{
		CompleteTasks:    s.stats.CompleteTasks,
		FailedTasks:      s.stats.FailedTasks,
		CancelTasks:      s.stats.FailedTasks,
		ActiveTasks:      s.stats.ActiveTasks,
		TotalSubmitTasks: s.stats.TotalSubmitTasks,
		QueueTasks:       s.stats.QueueTasks,
	}
}

func (s *Scheduler) work(id int) {
	defer func() {
		s.wg.Done()
		fmt.Printf("task:%d is stoping\n", id)
	}()

	fmt.Printf("task:%d is running\n", id)

	for {
		select {
		case <-s.ctx.Done():
			return

		case task, ok := <-s.TaskChan:
			if !ok {
				return
			}

			s.executeTask(task)
		}
	}
}

func (s *Scheduler) executeTask(task *Task) {
	taskResult := &TaskResult{
		TaskID:    task.TaskID,
		StartTime: time.Now(),
	}

	s.beforeExecuteTaks()

	ctx := s.ctx
	if task.Ctx != nil {
		mergeCtx, cancel := s.mergeCtx(task.Ctx, ctx)
		defer cancel()
		ctx = mergeCtx
	}
	if task.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(s.ctx, task.Timeout)
		defer cancel()
	}
	taskResult.Result, taskResult.Err = task.Fn(ctx)
	taskResult.EndTime = time.Now()
	taskResult.Duration = taskResult.EndTime.Sub(taskResult.StartTime)

	s.afterExecuteTaks(taskResult.Err)

	select {
	case s.ResultChan <- taskResult:
	case <-s.ctx.Done():
		return
	}
}

func (s *Scheduler) beforeSubmitTaks() {
	s.stats.mu.Lock()
	defer s.stats.mu.Unlock()

	s.stats.TotalSubmitTasks++
	s.stats.QueueTasks++
}

func (s *Scheduler) afterSubmitTaks() {
	s.stats.mu.Lock()
	defer s.stats.mu.Unlock()

	s.stats.QueueTasks--
}

func (s *Scheduler) beforeExecuteTaks() {
	s.stats.mu.Lock()
	defer s.stats.mu.Unlock()

	s.stats.ActiveTasks++
	s.stats.QueueTasks--
}

func (s *Scheduler) afterExecuteTaks(err error) {
	s.stats.mu.Lock()
	defer s.stats.mu.Unlock()

	s.stats.ActiveTasks--

	if err == nil {
		s.stats.CompleteTasks++
	} else if errors.Is(err, context.Canceled) {
		s.stats.CancelTasks++
	} else {
		s.stats.FailedTasks++
	}
}
