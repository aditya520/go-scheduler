package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

// TODO: Add more TODOs
// TODO: Solve all the TODOs

// Job ...
type Job func(ctx context.Context)

// Scheduler ...
type Scheduler struct {
	wg            *sync.WaitGroup
	cancellations []context.CancelFunc
}

// NewScheduler creates and return a New Scheduler object
func NewScheduler() *Scheduler {
	return &Scheduler{
		wg:            new(sync.WaitGroup),
		cancellations: make([]context.CancelFunc, 0),
	}
}

// Add starts goroutine which constantly calls provided job with interval delay
func (s *Scheduler) Add(ctx context.Context, j Job, interval time.Duration) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancellations = append(s.cancellations, cancel)

	s.wg.Add(1)
	go s.process(ctx, j, interval)
}

// Stop cancels all running jobs
func (s *Scheduler) Stop() {
	for _, cancel := range s.cancellations {
		cancel()
	}
	s.wg.Wait()
}

func (s *Scheduler) process(ctx context.Context, j Job, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			j(ctx)
		case <-ctx.Done():
			s.wg.Done()
			return
		}
	}
}

func main() {
	ctx := context.Background()

	worker := NewScheduler()
	worker.Add(ctx, task1, time.Second*5)
	worker.Add(ctx, task2, time.Second*10)

	quit := make(chan os.Signal, 1)
	// We can send multiple Interrupts. I think we should send interrupts = go routines spunned up.
	signal.Notify(quit, os.Interrupt)

	<-quit
	worker.Stop()
}

func task1(ctx context.Context) {
	// time.Sleep(time.Second * 1)
	fmt.Printf("Task 1 done %s\n", time.Now().String())
}

func task2(ctx context.Context) {
	// time.Sleep(time.Second * 1)
	fmt.Printf("Task 2 done %s\n", time.Now().String())
}
