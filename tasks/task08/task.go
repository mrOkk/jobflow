package task08

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-08--timeout-with-a-hidden-leak
// STATUS: READY FOR REVIEW
// BUG:
// Таймаут завершается раньше запроса в бд, а горутина с запросом об этом ничего не знает и возникает утечка
// FIX:
// Использовать контекст с таймаутом
// GUARANTEE:
// Контекст передается в горутину и отслеживает, когда запрос уже не нужен, в этом случае завершает работу горутины.
// В реальной системе вероятно контекст передавался бы и дальше непосредственно в функцию запроса
import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrTimeout = errors.New("query timeout")

// возможно здесь стоит возвращать ошибку, если ctx.Done.
// На реальной базе ошибку точно надо было бы обрабатывать, т.к. она при запросе в БД может вернуться
// а в сам запрос вероятно надо будет передать контекст
func queryDB(ctx context.Context) string {
	select {
	case <-ctx.Done():
		return ""
	case <-time.After(500 * time.Millisecond):
		return "result"
	}
}

func load() (string, error) {
	result := make(chan string)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(200)*time.Millisecond)
	defer cancel()

	go func() {
		value := queryDB(ctx)
		result <- value
	}()

	select {
	case value := <-result:
		return value, nil

	case <-ctx.Done():
		return "", ErrTimeout
	}
}

func Do() {
	value, err := load()

	fmt.Println("value:", value)
	fmt.Println("error:", err)
}
