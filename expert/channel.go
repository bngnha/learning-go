package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	ch := make(chan string, 5)

	strs := [5]string{"a", "b", "c", "d", "e"}

	for _, str := range strs {
		wg.Add(1)
		go sum(str, ch)
	}

	sum := 0
	for x := 0; x < 20; x++ {
		sum += x
	}
	for c := range ch {
		fmt.Println(c)
	}
	close(ch)
	wg.Wait()
	fmt.Println(sum)
}

func sum(str string, ch chan string) {
	defer wg.Done()
	t := 0
	for i := 0; i < 10; i++ {
		t += i
	}
	ch <- fmt.Sprintf("%s%d", str, t)
}
