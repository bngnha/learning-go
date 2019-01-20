package advanced

import (
	"fmt"
	"time"
)

func PlayWithVarDic() {
	fmt.Println(sum(2, 3, 5, 7, 11, 13, 17, 19, 23, 29))
	numbers := []int{31, 37, 41, 43, 47, 53, 59}
	fmt.Println(sum(numbers...))
}

func PlayWithCallback() {
	fmt.Printf("%v\n", timesTwo(func(i int) int {
		return i * 2
	}, 32))
}

func PlayWithClosure() {
	timeSince := initTimeSeq()
	fmt.Println(timeSince)

	time.Sleep(1 * time.Second)
	fmt.Println(timeSince())

	time.Sleep(120 * time.Millisecond)
	fmt.Println(timeSince())

	timeSince = initTimeSeq()
	time.Sleep(1300 * time.Millisecond)
	fmt.Println(timeSince())
}

func sum(numbers ...int) int {
	var total int
	for _, number := range numbers {
		total += number
	}
	return total
}

func timesTwo(f func(int) int, x int) int {
	return f(x * 2)
}

func initTimeSeq() func() int {
	t := time.Now().UnixNano()
	return func() int {
		return int(time.Now().UnixNano() - t)
	}
}
