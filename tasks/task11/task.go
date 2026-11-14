package task11

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-11--shutdown-means-finished
// STATUS: TODO
// BUG:
//
// FIX:
//
// GUARANTEE:
//
import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

type Service struct {
	jobs chan Job
	wg   sync.WaitGroup
}

func NewService(workers int) *Service {
	s := &Service{
		jobs: make(chan Job, 16),
	}

	s.wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(id int) {
			defer s.wg.Done()

			for job := range s.jobs {
				time.Sleep(100 * time.Millisecond)

				fmt.Printf(
					"worker %d processed %d\n",
					id,
					job.ID,
				)
			}
		}(i)
	}

	return s
}

func (s *Service) Submit(job Job) {
	s.jobs <- job
}

func (s *Service) Shutdown() {
	close(s.jobs)

	fmt.Println("service stopped")
}

func Do() {
	service := NewService(4)

	for i := 0; i < 20; i++ {
		service.Submit(Job{ID: i})
	}

	service.Shutdown()
}
