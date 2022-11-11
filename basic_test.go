package grabbag

import (
	"fmt"
	"testing"
	"time"
)

// printing to stdout causes weird scheduler patterns
func TestBasic(t *testing.T) {
	bag := []int{1, 2, 3, 4, 5}

	numWorkers := 4

	queue := make(chan int, 1)

	signal := make(chan struct{})

	fmt.Printf("bag: %#v\nnumWorkers: %d ", bag, numWorkers)

	// launch boss
	go boss(bag, queue, signal)

	for i := 0; i < numWorkers; i++ {
		go worker(i, queue)
	}

	time.Sleep(1 * time.Second)
	close(signal)
}

// worker: reads from queue
// quits if queue is closed
func worker[T any](id int, queue chan T) {
	for {
		num, ok := <-queue
		if !ok {
			break
		}

		fmt.Printf("id: %d\nnum: %v\n\n", id, num)
	}
}

// boss: writes to `queue` from `bag`
// closing `signal` closes `queue` and quits
func boss[T any](bag []T, queue chan T, signal chan struct{}) {
outer: // really exit out of loop, fixes weird issue of closing `queue` twice
	for i := 0; ; i++ {
		fmt.Printf("boss: %d\n\n", i)
		select {
		case <-signal:
			close(queue)
			break outer
		default:
			queue <- bag[i%(len(bag)-1)]
		}
	}
}
