package parallel

import (
	"errors"
	"sync"
	"sync/atomic"
)

type Runner interface {
	AddTask(TaskFunc, OnErrorFunc) (int, error)
	Run()
	Done()
	Cancel()
	Errors() map[int]error
}

type TaskFunc func(int) error

type OnErrorFunc func(error)

type task struct {
	run     TaskFunc
	onError OnErrorFunc
	num     uint32
}

type runner struct {
	tasks     chan *task
	taskCount uint32

	cancel      chan struct{}
	maxParallel int
	failFast    bool

	errors map[int]error
}

func NewRunner(maxParallel int, capacity uint, failFast bool) *runner {
	consumers := maxParallel
	if consumers < 1 {
		consumers = 1
	}
	if capacity < 1 {
		capacity = 1
	}
	r := &runner{
		tasks:       make(chan *task, capacity),
		cancel:      make(chan struct{}),
		maxParallel: consumers,
		failFast:    failFast,
	}
	r.errors = make(map[int]error)
	return r
}

func NewBounedRunner(maxParallel int, failFast bool) *runner {
	return NewRunner(maxParallel, 1, failFast)
}

func (r *runner) AddTask(t TaskFunc, errorHandler OnErrorFunc) (int, error) {
	return r.addTask(t, errorHandler)
}

func (r *runner) addTask(t TaskFunc, errorHandler OnErrorFunc) (int, error) {
	nextCount := atomic.AddUint32(&r.taskCount, 1)
	task := &task{run: t, num: nextCount - 1, onError: errorHandler}

	select {
	case <-r.cancel:
		return -1, errors.New("Runner stopped!")
	default:
		r.tasks <- task
		return int(task.num), nil
	}
}

func (r *runner) Run() {
	var wg sync.WaitGroup
	var m sync.Mutex
	var once sync.Once
	for i := 0; i < r.maxParallel; i++ {
		wg.Add(1)
		go func(threadId int) {
			defer wg.Done()
			for t := range r.tasks {
				e := t.run(threadId)
				if e != nil {
					if t.onError != nil {
						t.onError(e)
					}
					m.Lock()
					r.errors[int(t.num)] = e
					m.Unlock()
					if r.failFast {
						once.Do(r.Cancel)
						break
					}
				}
			}
		}(i)
	}
	wg.Wait()
}

func (r *runner) Done() {
	close(r.tasks)
}

func (r *runner) Cancel() {
	// No more adding tasks
	close(r.cancel)
	// Consume all tasks left
	for len(r.tasks) > 0 {
		<-r.tasks
	}
}

// Returns a map of errors keyed by the task number
func (r *runner) Errors() map[int]error {
	return r.errors
}
