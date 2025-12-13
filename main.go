package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("loading file error")
	}

	line := ""
	for {
		chunk := make([]byte, 8)

		n, err := file.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break
			}
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
			fmt.Printf("read: %s\n", line)
			line = ""
			line += parts[1]
		}
	}
}
