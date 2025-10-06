package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	Name     string
	Interval time.Duration
	Function func(ctx context.Context) error
	ctx      context.Context
	cancel   context.CancelFunc
	running  bool
}

type Scheduler struct {
	jobs   map[string]*Job
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func New() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		jobs:   make(map[string]*Job),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Scheduler) AddJob(name string, interval time.Duration, fn func(ctx context.Context) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[name]; exists {
		return fmt.Errorf("job with name '%s' already exists", name)
	}

	ctx, cancel := context.WithCancel(s.ctx)
	job := &Job{
		Name:     name,
		Interval: interval,
		Function: fn,
		ctx:      ctx,
		cancel:   cancel,
		running:  false,
	}

	s.jobs[name] = job
	fmt.Printf("Job '%s' added with interval: %s\n", name, interval)

	return nil
}

func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, job := range s.jobs {
		if !job.running {
			job.running = true
			s.wg.Add(1)
			go s.runJob(job)
		}
	}

	fmt.Printf("Scheduler started with %d jobs\n", len(s.jobs))
}

func (s *Scheduler) StartJob(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, exists := s.jobs[name]
	if !exists {
		return fmt.Errorf("job '%s' not found", name)
	}

	if job.running {
		return fmt.Errorf("job '%s' is already running", name)
	}

	// Eğer job daha önce durdurulmuşsa, yeni context oluştur
	if job.ctx.Err() != nil {
		ctx, cancel := context.WithCancel(s.ctx)
		job.ctx = ctx
		job.cancel = cancel
	}

	job.running = true
	s.wg.Add(1)
	go s.runJob(job)

	fmt.Printf("Job '%s' started\n", name)
	return nil
}

func (s *Scheduler) runJob(job *Job) {
	defer func() {
		s.mu.Lock()
		job.running = false
		s.mu.Unlock()
		s.wg.Done()
	}()

	ticker := time.NewTicker(job.Interval)
	defer ticker.Stop()

	fmt.Printf("Job '%s' is running every %s\n", job.Name, job.Interval)

	for {
		select {
		case <-job.ctx.Done():
			fmt.Printf("Job '%s' stopped\n", job.Name)
			return
		case <-ticker.C:
			fmt.Printf("Executing job: %s\n", job.Name)
			if err := job.Function(job.ctx); err != nil {
				fmt.Printf("Error executing job '%s': %v\n", job.Name, err)
			}
		}
	}
}

func (s *Scheduler) StopJob(name string) error {
	s.mu.RLock()
	job, exists := s.jobs[name]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("job '%s' not found", name)
	}

	if !job.running {
		return fmt.Errorf("job '%s' is not running", name)
	}

	job.cancel()
	fmt.Printf("Job '%s' stop signal sent\n", name)
	return nil
}

func (s *Scheduler) RemoveJob(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, exists := s.jobs[name]
	if !exists {
		return fmt.Errorf("job '%s' not found", name)
	}

	if job.running {
		job.cancel()
	}
	delete(s.jobs, name)
	fmt.Printf("Job '%s' removed\n", name)
	return nil
}

func (s *Scheduler) Stop() {
	fmt.Println("Stopping scheduler...")
	s.cancel()
	s.wg.Wait()
	fmt.Println("Scheduler stopped")
}

func (s *Scheduler) ListJobs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobs := make([]string, 0, len(s.jobs))
	for name := range s.jobs {
		jobs = append(jobs, name)
	}
	return jobs
}

func (s *Scheduler) GetJobInfo(name string) (string, time.Duration, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, exists := s.jobs[name]
	if !exists {
		return "", 0, fmt.Errorf("job '%s' not found", name)
	}

	return job.Name, job.Interval, nil
}

func (s *Scheduler) IsJobRunning(name string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, exists := s.jobs[name]
	if !exists {
		return false, fmt.Errorf("job '%s' not found", name)
	}

	return job.running, nil
}
