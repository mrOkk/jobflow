package task03

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-03--worker-pool-deadlock
// STATUS: TO REVIEW
// BUG:
// Вижу 2 проблемы (вторую осознал не сразу)
// 1. Чтение данных из канала по рэнджу, но при этом каналы никто не закрывает
// так что из этого цикла нет выхода
// 2. дедлок еще до чтения результатов
// мы записываем таски и читаем результаты последовательно в одной горутине
// каналы небуферизованные (это не особо имеет значение, просто проблема проявится максимально быстро)
// воркеров всего 5, а задач 100, 5 горутин читают задачи и пытаются записать результат
// который не может быть прочитан
//
// FIX:
// Точно нужно вынести запись задач и чтение результатов разнести по разным горутинам
// у меня родился такой вариант
// запись в tasks выделяем в отдельную горутину
// когда записали все нужные задачи, закрываем канал tasks
// далее создаем wg для отслеживания, когда все воркеры закончат обработку tasks
// можно в горутине записи tasks после закрытия канала tasks, можно в отдельной
// ожидаем завершения воркеров и закрываем results
//
// GUARANTEE:
// чтения и записи происходят в разных горутинах
// так что не могут полностью блокировать друг друга
// все канала, закрываются, так что обработка циклов завершается

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID int
}

func worker(
	id int,
	tasks <-chan Task,
	results chan<- string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for task := range tasks {
		fmt.Printf(
			"worker %d started task %d\n",
			id,
			task.ID,
		)

		time.Sleep(50 * time.Millisecond)

		results <- fmt.Sprintf(
			"worker %d completed task %d",
			id,
			task.ID,
		)
	}
}

func Do() {
	const numTasks = 100
	const numWorkers = 5

	tasks := make(chan Task)
	results := make(chan string)
	wg := &sync.WaitGroup{}
	wg.Add(numWorkers)

	for i := 1; i <= numWorkers; i++ {
		go worker(i, tasks, results, wg)
	}

	go func() {
		for i := 1; i <= numTasks; i++ {
			tasks <- Task{ID: i}
		}
		close(tasks)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}

	fmt.Println("all tasks completed")
}
