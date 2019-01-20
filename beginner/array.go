package beginner

import "fmt"

var strArray = []string{"a", "b", "c", "d", "f"}
var intArray = []int{1, 2, 3, 5, 8}
var mapOne = map[int]string{}
var mapTwo = map[string]interface{}{}

func WorkWithArray() {
	// do this five times
	for i := 0; i != 5; i++ {

		// print the $th value of the intarray and the strarray
		fmt.Println(intArray[i], "\t", strArray[i])

		mapOne[intArray[i]] = strArray[i]
		mapTwo[strArray[i]] = mapOne
	}
	fmt.Println(mapOne)
	fmt.Println(mapTwo)
}
