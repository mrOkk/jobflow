package task09

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-09--cancellation-doesnt-cancel
// STATUS: READY FOR REVIEW
// BUG:
// Контекст передается в вокрек, но не используется.
// А так как канал с результатами небуфферизируемый и запись происходит, а чтение нет,
// то зависает воркер, следовательно не вызывается wg.Done, следовательно зависает ожидание wg.
// Получаем дедлок
// FIX:
// Использовать контекст по назначению.
// Меняем for range на бесконечный for, чтение job через select, ожидание через select, запись через select.
// При чтении джобов проверяем, открыт ли доступны ли еще значения в канале на случай, если канал закрыли
// GUARANTEE:
// Любая ожидающая операция производится через select и если контекст отменили раньше, то результат не получим
import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

func worker(
	ctx context.Context,
	jobs <-chan Job,
	results chan<- int,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			time.Sleep(200 * time.Millisecond)

			select {
			case <-ctx.Done():
				return
			case <-time.After(200 * time.Millisecond):
				select {
				case <-ctx.Done():
					return
				case results <- job.ID:
				}
			}
		}
	}
}

func Do() {
	ctx, cancel := context.WithCancel(context.Background())

	jobs := make(chan Job)
	results := make(chan int)

	var wg sync.WaitGroup
	wg.Add(1)

	go worker(
		ctx,
		jobs,
		results,
		&wg,
	)

	go func() {
		jobs <- Job{ID: 1}
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	wg.Wait()

	fmt.Println("worker stopped")
}
