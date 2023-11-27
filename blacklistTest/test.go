package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("blacklist.txt")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var fileLines []string

	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}

	file.Close()

	fmt.Println("Enter your password: ")
	var password string
	fmt.Scanln(&password)

	var found = ""

	for _, line := range fileLines {
		if password == line {
			found = line
		}
	}

	if found != "" {
		fmt.Println("\nYour password matched with " + found + "\nPlease change it")
	} else {
		fmt.Printf("\nSuccess")
	}
}
