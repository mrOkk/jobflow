package task10

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-10--concurrent-repository
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
)

type Repository struct {
	values map[int]int
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
	r.values[key] = value
}

func (r *Repository) Get(
	key int,
) (int, bool) {
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
