package task15

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-15---rate-limiter-with-a-timer-leak
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
	"context"
	"fmt"
	"time"
)

type Request struct {
	ID int
}

func run(
	ctx context.Context,
	requests <-chan Request,
	rate int,
) {
	for request := range requests {
		<-time.Tick(
			time.Second / time.Duration(rate),
		)

		fmt.Println(
			"processed:",
			request.ID,
		)
	}
}

func Do() {
	ctx := context.Background()

	requests := make(chan Request)

	go run(
		ctx,
		requests,
		5,
	)

	for i := 0; i < 20; i++ {
		requests <- Request{ID: i}
	}

	close(requests)

	time.Sleep(time.Second)
}
