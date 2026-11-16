package task13

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-13---deadlock-transfer
// STATUS: READY FOR REVIEW
// BUG:
// G1 try lock A1 - success | G2 try lock A2 - success
// G1 try lock A2 - wait for unlock | G2 try lock A1 - wait for unlock
// На этом этапе происходит взаимная блокаировка
// FIX:
// Выстроить мьютексы в очередь
// GUARANTEE:
// Операция должна проходить при блокировку обоих счетов,
// так что можно просто выбрать правило, как организуем порядок
// и придерживаться его.
// В таком случае гарантируется, что Lock и Unlock корректно разблокируются
// COMMENT
// Сходу решение не придумал, так что искал, как это стоит делать
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
	mu1, mu2 := getMutexInOrder(from, to)
	mu1.Lock()
	defer mu1.Unlock()

	time.Sleep(10 * time.Millisecond)

	mu2.Lock()
	defer mu2.Unlock()

	from.Balance -= amount
	to.Balance += amount
}

func getMutexInOrder(
	acc1 *Account,
	acc2 *Account,
) (mu1 *sync.Mutex, mu2 *sync.Mutex) {
	if acc1.ID > acc2.ID {
		return &acc1.mu, &acc2.mu
	} else {
		return &acc2.mu, &acc1.mu
	}
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
