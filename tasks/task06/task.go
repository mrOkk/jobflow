package task06

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-06--stop-all-workers
// STATUS: READY FOR REVIEW
// BUG:
// Запись в канал дает возможность прочитать значение только в 1 горутине
// Соответственно, другие не обработают "событие" остановки
// FIX:
// Закрыть канал вместо записи в него
// GUARANTEE:
// закрытый канал обрабатывается в любом количестве читателей и
// если не сделать возврат, то такой кейс будет спасить бесконечно
// при желании можно обработать признак доступности канала value, open := <-chan

// Про то, чем отличается запись и закрытие по сути ответ был
// одну единицу данных из канала получит только один читатель
// информацию о том, что канал закрыт получат все
// Вопрос 2: физически закрыть канал может кто угодно, у кого есть возможность записи
// т.е. если передать канал без ограничение только на чтение, то его можно и закрыть
// но в данном случае есть защита и такой тактики стоит придерживать
import (
	"fmt"
	"sync"
	"time"
)

func worker(
	id int,
	stop <-chan struct{},
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-stop:
			fmt.Println("worker stopped:", id)
			return

		case <-time.After(100 * time.Millisecond):
			fmt.Println("worker working:", id)
		}
	}
}

func Do() {
	const workers = 5

	stop := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go worker(i, stop, &wg)
	}

	time.Sleep(300 * time.Millisecond)

	close(stop)

	wg.Wait()

	fmt.Println("all workers stopped")
}
