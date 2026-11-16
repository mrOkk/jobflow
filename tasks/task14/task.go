package task14

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-14---race-free-but-still-wrong
// STATUS: READY FOR REVIEW
// BUG:
// У нас как будто бы при установки Complete выполняется тяжелая операция.
// В процессе горутина передает управление, в этот момент pending
// уже уменьшен, а done еще не увеличен
//
// FIX:
// Мне пока пришло в голову только переработать структуру Stats:
// убрать pending, а сделать обычный целочисленный total,
// и в стэпшоте получать pending как total - done
// А из операций оставить только s.done.Add(1)
// 2 изменяемых поля pending и done избыточны для такой структуры
// В этом случае проверка теряет смысл.
//
// Либо еще вариант получать снэпшот и менять значения под единым мьютексом
//
// GUARANTEE:
// Получаем операцию только над одним полем, которая по настоящему атомарная
import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

type Stats struct {
	total int64
	done  atomic.Int64
}

func NewStats(total int64) *Stats {
	s := &Stats{}
	s.total = total

	return s
}

func (s *Stats) Complete() {
	runtime.Gosched()

	s.done.Add(1)
}

func (s *Stats) Snapshot() (
	pending int64,
	done int64,
) {
	done = s.done.Load()
	pending = s.total - done
	return
}

func Do() {
	const total = 1000

	stats := NewStats(total)

	var wg sync.WaitGroup
	wg.Add(total)

	for i := 0; i < total; i++ {
		go func() {
			defer wg.Done()

			stats.Complete()
		}()
	}

	for {
		pending, done := stats.Snapshot()

		if pending+done != total {
			fmt.Printf(
				"BROKEN SNAPSHOT: pending=%d done=%d\n",
				pending,
				done,
			)

			break
		}

		if done == total {
			break
		}
	}

	wg.Wait()
}
