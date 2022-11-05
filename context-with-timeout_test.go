package belajar_golang_context

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func CreateCounterExampleTimeout(group *sync.WaitGroup, ctx context.Context) chan int {
	defer group.Done()
	group.Add(1)
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
				time.Sleep(1 * time.Second) //simulasi slow response
			}
		}
	}()
	return destination
}

func TestContextWithTimeout(t *testing.T) {
	group := &sync.WaitGroup{}
	fmt.Println("Total Goroutine = ", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	destination := CreateCounterExampleTimeout(group, ctx)
	fmt.Println("Total Goroutine = ", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter ", n)
	}

	time.Sleep(2 * time.Second)
	group.Wait()

	fmt.Println("Total Akhir Goroutine = ", runtime.NumGoroutine())
}
