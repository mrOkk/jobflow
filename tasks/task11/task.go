package task11

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-11--shutdown-means-finished
// STATUS: IN PROGRESS
// BUG:
// Исходя из контракта, выделю 2 проблемных момента
// 1. Конкретно в этом код оно не используется, но описано в доп сценарии
// Попытка добавить работу после shutdown приведет к панике
// 2. Работа не завершится. Несколько задач хоть и отправлятся в обработку,
// но те, которые не успел туда попасть игнорируются полностью.
// А тех, которые взяты, main не дожидается
// FIX:
// Заготовка для ожидания завершения воркеров есть и канал закрывается,
// но этого мало.
// Добавлю еще одну wg, на ожидание завершения работ,
// Признак, что сервис остановлен
// и Mutex, чтобы безопасно читать и устанавливать признак.
// При вызове Shutdown устанавливаем признак stopped, чтобы не принимать новые задачи,
// дожидаемся завершения обработки текущих задач,
// закрываем канал,
// дожидаемся завершения воркеров.
// GUARANTEE:
// В такой схеме мы соблюдаем требуемые контракты, принятая работа будет выполнена,
// новая не добавится
import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrStopped = errors.New("service stopped")

type Job struct {
	ID int
}

type Service struct {
	jobs    chan Job
	wwg     sync.WaitGroup
	jwg     sync.WaitGroup
	mu      sync.Mutex
	stopped bool
}

func NewService(workers int) *Service {
	s := &Service{
		jobs: make(chan Job, 16),
	}

	s.wwg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(id int) {
			defer s.wwg.Done()

			for job := range s.jobs {
				time.Sleep(100 * time.Millisecond)

				fmt.Printf(
					"worker %d processed %d\n",
					id,
					job.ID,
				)

				s.jwg.Done()
			}
		}(i)
	}

	return s
}

func (s *Service) Submit(job Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopped {
		return ErrStopped
	}

	s.jwg.Add(1)
	s.jobs <- job
	return nil
}

func (s *Service) Shutdown() {
	s.mu.Lock()
	if s.stopped {
		// already stopped
		return
	}

	s.stopped = false
	s.mu.Unlock()

	s.jwg.Wait()
	close(s.jobs)
	s.wwg.Wait()

	fmt.Println("service stopped")
}

func Do() {
	service := NewService(4)

	for i := 0; i < 20; i++ {
		err := service.Submit(Job{ID: i})
		if err != nil {
			break
		}
	}

	service.Shutdown()
}
