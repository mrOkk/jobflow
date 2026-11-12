package task07

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-07--broken-fan-in
// STATUS: TO REVIEW
// BUG:
// Во время слияния канал out закрывается в каждом вызове CopyValues
// Это гарантированно приводит к панике
// FIX:
// Создать wait group для копирования значений и закрывать, когда дождемся окончания слияния
// GUARANTEE:
// wg обеспечит гарантию того, что канал закроется только после всех операций записи
// соответственно, никакой паники не возникнет
import (
	"fmt"
	"sync"
)

func source(start, end int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for i := start; i <= end; i++ {
			out <- i
		}
	}()

	return out
}

func merge(
	a <-chan int,
	b <-chan int,
) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		wg.Wait()
		close(out)
	}()

	copyValues := func(in <-chan int) {
		defer wg.Done()

		for value := range in {
			out <- value
		}
	}

	go copyValues(a)
	go copyValues(b)

	return out
}

func Do() {
	a := source(1, 5)
	b := source(100, 105)

	for value := range merge(a, b) {
		fmt.Println(value)
	}
}
