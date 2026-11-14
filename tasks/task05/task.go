package task05

// https://github.com/xakepp35/go_backend_combat/blob/main/combat/03_CONCURRENCY.md#task-05--pipeline-that-never-finishes
// STATUS: READY FOR REVIEW
// BUG:
// Не закрываются каналы, соответственно нет выхода из циклов, горутины зависают
// FIX:
// Закрывать каналы после записи всего необходимого в них
// GUARANTEE:
// Закрытие каналов происходит по цепочке стадий, range завершается
// горутины выходят, пайплайн не зависает
import "fmt"

func generate(n int) <-chan int {
	out := make(chan int)

	go func() {
		for i := 1; i <= n; i++ {
			out <- i
		}
		close(out)
	}()

	return out
}

func even(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for value := range in {
			if value%2 == 0 {
				out <- value
			}
		}
		close(out)
	}()

	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for value := range in {
			out <- value * value
		}
		close(out)
	}()

	return out
}

func Do() {
	numbers := generate(10)
	evens := even(numbers)
	squares := square(evens)

	for value := range squares {
		fmt.Println(value)
	}

	fmt.Println("pipeline completed")
}
