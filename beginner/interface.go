package beginner

import "fmt"

func FunctionWithInterface() {
	var myInterface = [3]interface{}{}

	myInterface[0] = 23
	myInterface[1] = "test"
	myInterface[2] = false
	fmt.Printf("Data %v\n", myInterface)

	for _, v := range myInterface {
		printInterfaceData(v)
	}
}

func printInterfaceData(myInterface interface{}) {
	switch t := myInterface.(type) {
	case string:
		fmt.Print("Type: string\t")
	case int:
		fmt.Print("Type: int\t")
	case bool:
		fmt.Print("Type: bool\t")
	default:
		fmt.Printf("Type: %v\t", t)
	}

	fmt.Printf("Data: %#v\n", myInterface)
}
