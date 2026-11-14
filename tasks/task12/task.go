package task12

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-12---initialize-exactly-once
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

type Client struct {
	ID int
}

var client *Client

func getClient() *Client {
	if client == nil {
		time.Sleep(100 * time.Millisecond)

		client = &Client{
			ID: 42,
		}
	}

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
