package task04

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-04--fake-backpressure
// STATUS: TO REVIEW
// BUG:
// Опасная конструкция: запись в канал из горутины, а сразу после закрытие
// Здесь скорее всего получим панику записи в закрытый канал
// еще в условии говорится Исправьте программу так, чтобы количество входных Job не превращалось в количество созданных goroutines
// это будет второй частью фикса
//
// FIX:
// Добавить wait group на запись в канал и закрывать канал только после окончания
// лучше в отдельной горутине
//
// Для второй части я думаю добавить семафор
// Как понимаю, для go наиболее идеоматичный вариант - это буфферезированный канал
// Сделаем его на тоже число элементов, что и размер очереди работ
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

func worker(
	jobs <-chan Job,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for job := range jobs {
		time.Sleep(10 * time.Millisecond)
		fmt.Println("processed:", job.ID)
	}
}

func Do() {
	const workers = 4
	const queueSize = 8
	const jobsCount = 1_000

	jobs := make(chan Job, queueSize)
	sem := make(chan struct{}, queueSize)

	var wg sync.WaitGroup
	wg.Add(workers)
	var jwg sync.WaitGroup
	jwg.Add(jobsCount)

	for i := 0; i < workers; i++ {
		go worker(jobs, &wg)
	}

	for i := 0; i < jobsCount; i++ {
		job := Job{ID: i}
		sem <- struct{}{}

		go func() {
			defer jwg.Done()
			jobs <- job
			<-sem
		}()
	}

	close(sem)

	go func() {
		jwg.Wait()
		close(jobs)
	}()

	wg.Wait()
}
