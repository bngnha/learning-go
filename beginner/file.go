package main

import (
	"os"
)

// Write your content to file
func writeToFile(path string, content string) {
	fp := createFile(path)

	fp.WriteString(content)

	fp.Close()
}

func createFile(path string) *os.File {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fp, err := os.Create(path)
		if err != nil {
			panic(err)
		}

		return fp
	}

	// open file in write mode
	fp, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	return fp
}

func main() {
	writeToFile()
}
