package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		os.Exit(1)
	}
	linesCh := getLinesChannel(file)

	for line := range linesCh {
		fmt.Println("read:", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		buffer := make([]byte, 8)
		currentLine := ""

		for {
			n, err := f.Read(buffer)
			if err != nil {
				if err == io.EOF {
					if currentLine != "" {
						ch <- currentLine
					}
					return
				}
				fmt.Println(err)
				return
			}

			currentString := string(buffer[:n])
			if strings.Contains(currentString, "\n") {
				splitString := strings.Split(currentString, "\n")
				ch <- currentLine + splitString[0]
				currentLine = splitString[1]
			} else {
				currentLine = currentLine + currentString
			}
		}
	}()
	return ch
}
