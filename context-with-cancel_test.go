package belajar_golang_context

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func CreateCounter(group *sync.WaitGroup, ctx context.Context) chan int {
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
			}
		}
	}()
	return destination
}

func TestContextWithCancel(t *testing.T) {
	group := &sync.WaitGroup{}
	fmt.Println("Total Goroutine = ", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(group, ctx)
	fmt.Println("Total Goroutine = ", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter ", n)
		if n == 10 {
			break
		}
	}

	cancel()
	time.Sleep(2 * time.Second)
	group.Wait()

	fmt.Println("Total Akhir Goroutine = ", runtime.NumGoroutine())
}
