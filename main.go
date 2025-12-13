package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer f.Close() // always close file when goroutine ends
		defer close(ch) // always close channel when done
		line := ""
		for {
			chunk := make([]byte, 8)
			n, err := f.Read(chunk)

			if err != nil {
				if err == io.EOF {
					// f.Close()
					// close(ch)
					break
				}
				// f.Close()
				// close(ch)
				break
			}
			str := string(chunk[:n])
			parts := strings.Split(str, "\n")
			if len(parts) == 1 {
				line += parts[0]
				continue
			}
			if len(parts) == 2 {
				line += parts[0]
				ch <- line
				line = ""
				line += parts[1]
			}
		}
	}()
	return ch
}

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("loading file error")
	}
	ch := getLinesChannel(file)
	for line := range ch {
		fmt.Printf("read: %s\n", line)
	}
}
