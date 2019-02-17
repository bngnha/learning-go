package main

import "fmt"

type stack struct {
	nodes []string
	count int
}

func (s *stack) push(n string) {
	s.nodes = append(s.nodes[:s.count], n)
	s.count++
}

func (s *stack) pop() string {
	if s.count == 0 {
		return ""
	}
	s.count--
	return s.nodes[s.count]
}

func lifo() *stack {
	return &stack{}
}

func main() {
	lf := lifo()
	lf.push("Hello")
	lf.push("World")
	lf.push("From")
	lf.push("Golang")

	fmt.Println(lf.pop())
	fmt.Println(lf.pop())
}
