package expert

import (
	"fmt"
	"sync"
	"time"
)

var start time.Time

func init() {
	start = time.Now()
}

// PlayWithChannel function
func PlayWithChannel() {
	fmt.Println("[main] started", time.Since(start))

	chan1 := make(chan string)
	chan2 := make(chan string)

	go service1(chan1)
	go service2(chan2)

	select {
	case res := <-chan1:
		fmt.Println("Response from service 1:", res, time.Since(start))
	case res := <-chan2:
		fmt.Println("Respnse from service 2:", res, time.Since(start))
	}

	fmt.Println("main() stopped!")
}

func service1(c chan string) {
	time.Sleep(3 * time.Second)
	c <- "Hello from service 1"
}

func service2(c chan string) {
	time.Sleep(5 * time.Second)
	c <- "Hello from service 2"
}

func square(c chan int) {
	fmt.Println("[square] reading")

	num := <-c
	c <- num * num
}

func cube(c chan int) {
	fmt.Println("[cube] reading")

	num := <-c
	c <- num * num * num
}

func greet(cc chan chan string) {
	c := make(chan string)
	cc <- c
}

func greeter(c chan string) {
	fmt.Println("Hello ", <-c)
}

// ChannelFunc function
func ChannelFunc() {

	var wg sync.WaitGroup

	ch := make(chan string, 5)

	strs := [5]string{"a", "b", "c", "d", "e"}

	for _, str := range strs {
		wg.Add(1)
		go sum(str, ch, &wg)
	}

	total := 0
	for x := 0; x < 20; x++ {
		total += x
	}
	for c := range ch {
		fmt.Println(c)
	}
	close(ch)
	wg.Wait()
	fmt.Println(total)
}

func sum(str string, ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	t := 0
	for i := 0; i < 10; i++ {
		t += i
	}
	ch <- fmt.Sprintf("%s%d", str, t)
}
