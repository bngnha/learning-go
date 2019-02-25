package expert

import (
	"fmt"
	"sync"
)

// PlayWithChannel function
func PlayWithChannel() {
	c := make(chan string, 1)
	//c1 := make(chan string, 1)
	c2 := make(chan string, 1)
	//var c chan int

	//fmt.Printf("type of `c` is %T\n", c)
	//fmt.Printf("value of `c` is %v\n", c)

	//for i := 0; i < 2; i++ {
	go greet(c, c2)
	//}
	c <- "C"
	//c1 <- "John"
	//c1 <- "Mary"

	//c2 <- "Test"

	//time.Sleep(time.Second * 1)
	//fmt.Println("First Hello " + <-c)
	fmt.Println("main() stopped!")
}

func greet(c1, c2 chan string) {
	//for x := range c {
	//	fmt.Println("c")
	//	fmt.Println("Hello " + x + "!")
	//}
	//fmt.Println("d")
	for {
		select {
		case x := <-c1:
			fmt.Println(x)
			//close(c)
			//<-c
		case y := <-c2:
			fmt.Println(y)
			//c2 <- y
		}
	}
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
