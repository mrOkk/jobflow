package task10

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-10--concurrent-repository
// STATUS: READY FOR REVIEW
// BUG:
// map не является потокобезопасным типом, так что такое использование приведет к чему-то нехорошему
// FIX:
// Добавить в репозиторий rw mutex и использовать его
// GUARANTEE:
// WRMutex позволяет читать любому количеству читателей и блокировать все только при записи,
// что помогает избежать fatal error
import (
	"fmt"
	"sync"
)

type Repository struct {
	values map[int]int
	rw     sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		values: make(map[int]int),
	}
}

func (r *Repository) Set(
	key int,
	value int,
) {
	r.rw.Lock()
	defer r.rw.Unlock()

	r.values[key] = value
}

func (r *Repository) Get(
	key int,
) (int, bool) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	value, ok := r.values[key]
	return value, ok
}

func Do() {
	repo := NewRepository()

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(2)

		go func(id int) {
			defer wg.Done()

			repo.Set(id, id*10)
		}(i)

		go func(id int) {
			defer wg.Done()

			repo.Get(id)
		}(i)
	}

	wg.Wait()

	fmt.Println("completed")
}
