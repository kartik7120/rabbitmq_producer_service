package tests

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_sample(t *testing.T) {

	t.Run("Understanding how does select block works", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Simulate some work in a goroutine that takes 3 seconds
		resultChan := make(chan string, 1)

		go func() {
			// Simulate a slow task
			time.Sleep(3 * time.Second)
			resultChan <- "Task completed"
		}()

		// Use select to wait for either the task to finish or the context to timeout
		select {
		case res := <-resultChan:
			fmt.Println(res)
		case <-ctx.Done():
			fmt.Println("Context timeout:", ctx.Err()) // will print: context deadline exceeded
		default:
			fmt.Println("I will always run")
		}

	})
}
