package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var (
	firstName, lastName, s1 string
	i1                      int
	f3                      float32
	input                   = "56.12 / 5212 / Go"
	format                  = "%f / %d / %s"
)

func main() {
	fmt.Println("Please enter your full name: ")
	fmt.Scanln(&firstName, &lastName)
	fmt.Printf("Hi %s %s!\n", firstName, lastName) // Hi Chris Naegels
	fmt.Sscanf(input, format, &f3, &i1, &s1)
	fmt.Println("From the string we read: ", f3, i1, s1)

}

// https://medium.com/swlh/the-most-common-mistakes-with-read-file-in-golang-be87239fd03b
/**
how to extract read file from Go, easy to be test, using Reader as parameter
*/

func count(reader *bufio.Reader) (int, error) {
	count := 0
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			switch err {
			default:
				return 0, fmt.Errorf("unable to read", err)
			case io.EOF:
				return count, nil
			}
		}
		if len(line) == 0 {
			count++
		}
	}
}

func ReadFile() error {
	filename := os.Getenv("fileExample")
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("unable to open %s", filename, err)
	}
	defer file.Close()
	count, err := count(bufio.NewReader(file))
	println("count=", count)
	return nil
}

/*
Please enter your full name:
Chen zhe
Hi Chen zhe!
From the string we read:  56.12 5212 Go
*/
