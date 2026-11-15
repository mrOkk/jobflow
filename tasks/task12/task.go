package task12

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-12---initialize-exactly-once
// STATUS: READY FOR REVIEW
// BUG:
// На получение клиента нет никакой синхронизакии, при этом операция тяжелая,
// возникает гонка данных, данные много раз перетираются,
// один вызов может получить отличное от другого значения
// FIX:
// sync.Once для создания экземпляра вместо проверки на null
// GUARANTEE:
// sync.Once гарантированно выполняется только один раз, так что мы ни
// при каких обстаятельствах не можем получить отличающиеся значения
import (
	"fmt"
	"sync"
	"time"
)

type Client struct {
	ID int
}

var client *Client
var once sync.Once

func getClient() *Client {
	once.Do(func() {
		time.Sleep(100 * time.Millisecond)

		client = &Client{
			ID: 42,
		}
	})

	return client
}

func Do() {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			c := getClient()
			fmt.Printf("%p\n", c)
		}()
	}

	wg.Wait()
}
