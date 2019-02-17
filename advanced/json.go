package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	var strings []string
	var jsonString = `["hello", "world", "from", "golang"]`

	err := json.Unmarshal([]byte(jsonString), &strings)

	if err != nil {
		fmt.Println("error while unmarshal")
		os.Exit(2)
	}
	fmt.Println(strings)

	//convert object to byte-string
	jsonData, err := json.Marshal(strings)

	if err != nil {
		fmt.Println("Error while marshal")
		os.Exit(2)
	}
	fmt.Println(jsonData)

	err = json.Unmarshal([]byte(jsonData), &strings)
	if err != nil {
		fmt.Println("Error while unmarshal again")
		os.Exit(2)
	}
	fmt.Println(strings)
}
