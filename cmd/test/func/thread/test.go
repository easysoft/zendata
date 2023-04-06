package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	total = 251

	size = 100
	c    = SafeCounter{v: make(map[string]int)}
)

func main() {
	// runtime.GOMAXPROCS(1)

	fmt.Println("=== start ===")

	threadCount := getTheadCount(total, size)

	var waitGroup sync.WaitGroup
	waitGroup.Add(threadCount)

	for i := 0; i < threadCount; i++ {
		index := i
		c.Inc("count")

		fmt.Printf("index = %d\n", index)

		recordNumb := getTheadRecordNum(total, size, threadCount, index)

		go func() {
			defer func() {
				c.Dec("count")
				waitGroup.Done()
			}()

			for recordIndex := 0; recordIndex < recordNumb; recordIndex++ {
				fmt.Printf("gen %d - %d\n", index, recordIndex+1)
				time.Sleep(time.Millisecond * 100)
			}
		}()
	}

	go func() {
		for {
			fmt.Printf("=thread count= %d\n", c.Value("count"))
			time.Sleep(time.Millisecond * 50)
		}
	}()

	waitGroup.Wait()
	fmt.Println("=== end ===")
}

func getTheadCount(total, size int) (threadCount int) {
	threadCount = total / size
	if total%size > 0 {
		threadCount++
	}

	return
}
func getTheadRecordNum(total, size, theadCount, index int) (num int) {
	if index == theadCount-1 {
		num = total % size
	} else {
		num = size
	}

	return
}

type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	c.v[key]++
	c.mu.Unlock()
}
func (c *SafeCounter) Dec(key string) {
	c.mu.Lock()
	c.v[key]--
	c.mu.Unlock()
}
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v[key]
}
