package expert

import (
	"fmt"
	"time"
)

// FanIn function
func FanIn() {
	c := make(chan int)
	o := make(chan int)

	go producer(c, 100*time.Millisecond)
	go producer(c, 150*time.Millisecond)
	go reader(o)

	for x := range c {
		o <- x
	}

	<-time.After(5 * time.Second)

}

func producer(c chan int, d time.Duration) {
	var i int
	for {
		c <- i
		i++
		time.Sleep(d)
	}
}

func reader(out chan int) {
	for o := range out {
		fmt.Println(o)
	}
}
