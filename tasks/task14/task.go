package task14

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-14---race-free-but-still-wrong
// STATUS: TODO
// BUG:
//
//
// FIX:
//
//
// GUARANTEE:
//
import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

type Stats struct {
	pending atomic.Int64
	done    atomic.Int64
}

func NewStats(total int64) *Stats {
	s := &Stats{}
	s.pending.Store(total)

	return s
}

func (s *Stats) Complete() {
	s.pending.Add(-1)

	runtime.Gosched()

	s.done.Add(1)
}

func (s *Stats) Snapshot() (
	pending int64,
	done int64,
) {
	return s.pending.Load(), s.done.Load()
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
