package grabbag

import (
	"fmt"
	"testing"
	"time"
)

func TestGenuis(t *testing.T) {
	bag := []int{1, 2, 3, 4, 5}

	numWorkers := 5

	queue := make(chan int, 1)

	signal := make(chan struct{})

	fmt.Printf("bag: %#v\nnumWorkers: %d ", bag, numWorkers)

	// launch boss
	go boss(bag, queue, signal)

	for i := 0; i < numWorkers; i++ {
		go genuis(i, queue)
	}

	time.Sleep(1 * time.Second)
	close(signal)
}

// genuis: reads from `queue` twice
// quits if `queue` is closed
func genuis[T any](id int, queue chan T) {
	for {
		num1, ok := <-queue
		if !ok {
			break
		}

		num2, ok := <-queue

		fmt.Printf("id: %d\nnum1: %v\nnum2: %v\n\n", id, num1, num2)
	}
}

// boss: writes to `queue` from `bag`
// closing `signal` closes `queue` and quits
func boss[T any](bag []T, queue chan T, signal chan struct{}) {
outer: // really exit out of loop, fixes weird issue of closing `queue` twice
	for i := 0; ; i += 2 {
		fmt.Printf("boss: %d\n\n", i)
		select {
		case <-signal:
			close(queue)
			break outer
		default:
			queue <- bag[i%(len(bag)-1)]

			// to fix the worker having 2 of the same element
			i++
			queue <- bag[i%(len(bag)-1)]
		}
	}
}
