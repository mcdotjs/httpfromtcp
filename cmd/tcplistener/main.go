package main

import (
	"fmt"
	"io"
	"log"
	"mirectm/httpfromtcp/internal/request"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("listening error")
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}
		fmt.Println("connection has been accepted from", conn.RemoteAddr())

		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("some req line error")
		}
		fmt.Println("Request line:")
		fmt.Println("- Method:", r.RequestLine.Method)
		fmt.Println("- Target:", r.RequestLine.RequestTarget)
		fmt.Println("- Version:", r.RequestLine.HttpVersion)

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

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
				if line != "" {
					ch <- line
				}
				if err == io.EOF {
					// f.Close()
					// close(ch)
					break
				}
				// f.Close()
				// close(ch)
				//break
				return
			}
			//n je vaccinou 8, ale hocikedy to bude menej (alebo error)
			str := string(chunk[:n])
			//tu musim checknut kde konci lina
			parts := strings.Split(str, "\n")

			for i := 0; i < len(parts)-1; i++ {
				//a to poslat chanellom vonku
				ch <- fmt.Sprintf("%s%s", line, parts[i])
				line = ""
			}
			// a prepandnut poslednu cast k dalsej line
			line += parts[len(parts)-1]
			// if len(parts) == 1 {
			// 	line += parts[0]
			// 	continue
			// }
			// if len(parts) == 2 {
			// 	line += parts[0]
			// }
		}
	}()
	return ch
}
