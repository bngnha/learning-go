package expert

import (
	"fmt"
	"time"
)

// PingPong function
func PingPong() {
	table := make(chan *ball)

	go player("ping", table)
	go player("pong", table)

	table <- new(ball)

	time.Sleep(1 * time.Second)
	<-table

	panic("Show me the stacks")
}

type ball struct{ hits int }

func player(name string, table chan *ball) {
	for {
		ball := <-table
		ball.hits++
		fmt.Println(name, ball.hits)
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}
