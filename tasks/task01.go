package tasks

import (
	"fmt"
	"sync"
)

// BUG:
// нет гарантии корректного подсчета из-за гонки данных
//
// FIX:
// использовать атомик или мьютекс
// атомик для каждого инкремента может будет накладным для процессора
// так что либо пересчитать математике и использовать CompareAndSwap
// или использовать mutex, тогда логику менять не требуется
//
// GUARANTEE:
// В случае с мьютексом мы по сути выстраиваем очередь выполнения
func Task01() {
	const goroutines = 1000
	const increments = 1000

	var wg sync.WaitGroup
	var counter int64
	var mu sync.Mutex

	for i := 0; i < goroutines; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			mu.Lock()
			defer mu.Unlock()

			for j := 0; j < increments; j++ {
				counter++
			}
		}()
	}

	wg.Wait()

	fmt.Println("counter:", counter)
}
