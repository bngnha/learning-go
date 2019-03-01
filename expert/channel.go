package expert

import (
	"fmt"
	"sync"
	"time"
)

var start time.Time
var i int

func init() {
	start = time.Now()
}

// PlayWithChannel function
func PlayWithChannel() {
	fmt.Println("[main] started", time.Since(start))

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go worker(&wg)
	}
	wg.Wait()

	fmt.Println("main() stopped!")
}

func worker(wg *sync.WaitGroup) {
	i++
	wg.Done()
}

func service1(c chan string) {
	//time.Sleep(3 * time.Second)
	fmt.Println("Service 1 started")
	c <- "Hello from service 1"
}

func service2(c chan string) {
	//time.Sleep(5 * time.Second)
	fmt.Println("Service 2 started")
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
