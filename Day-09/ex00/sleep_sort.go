package ex00

import (
	"sync"
	"time"
)

func sleepSort(nums []int) chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(nums))
	for _, num := range nums {
		go func(n int) {
			defer wg.Done()
			time.Sleep(time.Duration(n) * time.Second)
			out <- n
		}(num)
	}
	go func(){
		wg.Wait()
		close(out)
	}()	
	return out
}
