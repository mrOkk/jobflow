package task02

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-02--waitgroup-never-ends
// STATUS: READY FOR REVIEW
// BUG:
// waitGroup передается по значению
// соответственно внутренний счетчик уменьшается не у того инстанса, у которого увеличивается
// а для оригинального инстанса не меняется. Соответственно wg.Wait никогда не закончится
// небольшая проблема, что любая IDE подсвечивает этот баг
//
// FIX:
// 1. использовать указатель на wg
// 2. вызывать не сразу fetch с передачей wg, а
// go func(){
// defer wg.Done()
// fetch(url, results)
// }()
// или вооще обертку wg.Go, только тогда надо будет убрать ручной wg.Add
//
// GUARANTEE:
// указатель указывает на одну область в памяти
// так что если сходится добавляемое количество с количеством вызовов wg.Done()
// то все отработает как надо
import (
	"fmt"
	"sync"
	"time"
)

type Result struct {
	URL    string
	Status int
}

func fetch(
	url string,
	wg *sync.WaitGroup,
	results chan<- Result,
) {
	defer wg.Done()

	time.Sleep(100 * time.Millisecond)

	results <- Result{
		URL:    url,
		Status: 200,
	}
}

func Do() {
	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://golang.org",
	}

	results := make(chan Result, len(urls))

	wg := &sync.WaitGroup{}
	wg.Add(len(urls))

	for _, url := range urls {
		//wg.Go(func() { fetch(url, results) })
		go fetch(url, wg, results)
	}

	wg.Wait()
	close(results)

	for result := range results {
		fmt.Printf(
			"%s -> %d\n",
			result.URL,
			result.Status,
		)
	}

	fmt.Println("done")
}
