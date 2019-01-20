package advanced

import (
	"fmt"
	"math"
)

func PlayWithOOPBasic() {
	r := rectangle{width: 2, height: 3}
	s := square{length: 3}
	c := circle{radius: 4}
	//q := cuboid{width: 3, height: 2, length: 4}

	geocalc(r)
	geocalc(s)
	geocalc(c)
}

type geo interface {
	area() float64
	extent() float64
	volume() float64
}
type rectangle struct {
	width, height float64
}

type square struct {
	length float64
}

type circle struct {
	radius float64
}

type cuboid struct {
	width, height, length float64
}

func (r rectangle) area() float64 {
	return r.width * r.height
}

func (r rectangle) extent() float64 {
	return 2*r.width + 2*r.height
}

func (r rectangle) volume() float64 {
	return 0
}

func (s square) area() float64 {
	return s.length * s.length
}

func (s square) extent() float64 {
	return 4 * s.length
}

func (s square) volume() float64 {
	return 0
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c circle) extent() float64 {
	return math.Pi * (c.radius + c.radius)
}

func (c circle) volume() float64 {
	return 0
}

func geocalc(g geo) {
	fmt.Printf("%#v\t%#v\t%#v\t%#v\n", g, g.area(), g.extent(), g.volume())
}
