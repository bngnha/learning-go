package advanced

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
