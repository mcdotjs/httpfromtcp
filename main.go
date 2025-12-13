package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	//load once
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("loading file error")
	}

	for {
		chunk := make([]byte, 8)

		n, err := file.Read(chunk)
		if err != nil {
			if err == io.EOF{
				break
			}
			break
		}
		// fmt.Printf("read: %s\n", string(chunk))     // "Hello!\x00\x00" ❌
		// fmt.Printf("read: %s\n", string(chunk[:n])) // "Hello!"
		fmt.Printf("read: %s\n", string(chunk[:n]))
	}
}
