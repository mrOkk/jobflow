package task08

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-08--timeout-with-a-hidden-leak
// STATUS: TODO
// BUG:
//
// FIX:
//
// GUARANTEE:
//
import (
	"errors"
	"fmt"
	"time"
)

var ErrTimeout = errors.New("query timeout")

func queryDB() string {
	time.Sleep(500 * time.Millisecond)
	return "result"
}

func load() (string, error) {
	result := make(chan string)

	go func() {
		value := queryDB()
		result <- value
	}()

	select {
	case value := <-result:
		return value, nil

	case <-time.After(200 * time.Millisecond):
		return "", ErrTimeout
	}
}

func Do() {
	value, err := load()

	fmt.Println("value:", value)
	fmt.Println("error:", err)
}
