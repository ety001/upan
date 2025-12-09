package tasks

import (
	"context"
	"log"
	"sync"
	"time"
)

// Task represents a scheduled task
type Task interface {
	Name() string
	Run(ctx context.Context) error
	Interval() time.Duration
}

// Scheduler manages scheduled tasks
type Scheduler struct {
	tasks   []Task
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	started bool
	mu      sync.Mutex
}

// NewScheduler creates a new task scheduler
func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		tasks:  make([]Task, 0),
		ctx:    ctx,
		cancel: cancel,
	}
}

// Register adds a task to the scheduler
func (s *Scheduler) Register(task Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = append(s.tasks, task)
	log.Printf("Registered task: %s (interval: %v)", task.Name(), task.Interval())
}

// Start begins running all registered tasks
func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		log.Println("Scheduler is already started")
		return
	}

	s.started = true
	log.Printf("Starting scheduler with %d tasks", len(s.tasks))

	for _, task := range s.tasks {
		s.wg.Add(1)
		go s.runTask(task)
	}

	log.Println("Scheduler started successfully")
}

// Stop stops all running tasks
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		return
	}

	log.Println("Stopping scheduler...")
	s.cancel()
	s.wg.Wait()
	s.started = false
	log.Println("Scheduler stopped")
}

// runTask runs a single task in a loop
func (s *Scheduler) runTask(task Task) {
	defer s.wg.Done()

	// Run immediately on start
	if err := s.executeTask(task); err != nil {
		log.Printf("Error executing task %s: %v", task.Name(), err)
	}

	// Then run on interval
	ticker := time.NewTicker(task.Interval())
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			log.Printf("Task %s stopped", task.Name())
			return
		case <-ticker.C:
			if err := s.executeTask(task); err != nil {
				log.Printf("Error executing task %s: %v", task.Name(), err)
			}
		}
	}
}

// executeTask executes a task with error handling and logging
func (s *Scheduler) executeTask(task Task) error {
	start := time.Now()
	log.Printf("Executing task: %s", task.Name())

	err := task.Run(s.ctx)
	duration := time.Since(start)

	if err != nil {
		log.Printf("Task %s failed after %v: %v", task.Name(), duration, err)
		return err
	}

	log.Printf("Task %s completed successfully in %v", task.Name(), duration)
	return nil
}

// GetTaskCount returns the number of registered tasks
func (s *Scheduler) GetTaskCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.tasks)
}

// IsRunning returns whether the scheduler is currently running
func (s *Scheduler) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.started
}

