package task13

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-13---deadlock-transfer
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

type Account struct {
	ID      int
	mu      sync.Mutex
	Balance int
}

func transfer(
	from *Account,
	to *Account,
	amount int,
) {
	from.mu.Lock()
	defer from.mu.Unlock()

	time.Sleep(10 * time.Millisecond)

	to.mu.Lock()
	defer to.mu.Unlock()

	from.Balance -= amount
	to.Balance += amount
}

func Do() {
	a := &Account{
		ID:      1,
		Balance: 1000,
	}

	b := &Account{
		ID:      2,
		Balance: 1000,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		transfer(a, b, 100)
	}()

	go func() {
		defer wg.Done()
		transfer(b, a, 50)
	}()

	wg.Wait()

	fmt.Println(
		a.Balance,
		b.Balance,
	)
}
