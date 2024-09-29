package main

import (
	"fmt"
	"sync"
)

var (
	count int
)

func incrementWithoutMutex() {
	count++
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			incrementWithoutMutex()
		}()
	}

	wg.Wait()
	fmt.Println("Final count:", count) // Might not be 10 due to race conditions
}
