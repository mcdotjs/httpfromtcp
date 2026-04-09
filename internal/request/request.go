package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	ParserState int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const CRLF string = "\r\n"

const (
	initialized = iota
	done
)

func ParseRequestLine(b []byte) (int, error) {
	idxOfRegisteredNurse := bytes.Index(b, []byte(CRLF))
	if idxOfRegisteredNurse == -1 {
		//NOTE: ak nenajde crlf .... just give me more data :)
		return 0, nil
	}

	return idxOfRegisteredNurse, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{
		ParserState: initialized,
	}
	bufferSize := 8
	readToIndex := 0
	//parsedBytes1 := 0
	buffer := make([]byte, bufferSize)
	for r.ParserState != done {

		if len(buffer) == readToIndex {
			fmt.Printf("buffer growing: %q\n", buffer)
			bufferSize *= 2
			temp := make([]byte, bufferSize)
			copy(temp, buffer)
			buffer = temp
			fmt.Printf("buffer growing1: %q\n", buffer)
		}
		readBytes, err := reader.Read(buffer[readToIndex:])

		fmt.Printf("readBytes: %d\n", readBytes)
		if err != nil {
			if err == io.EOF {
				r.ParserState = done
				break
			}
			return nil, fmt.Errorf("error reading to buffer")
		}

		fmt.Printf("buffer: %q\n", buffer)
		readToIndex += readBytes

		parsedBytes, err := ParseRequestLine(buffer)

		if parsedBytes != 0 {
			r.ParserState = done

			fmt.Println("out:", readToIndex, parsedBytes)
			reqLine, err := ParseRequestLineFromString(string(buffer[:parsedBytes]))
			fmt.Println("RESULT:", reqLine)
			if err != nil {
				return nil, fmt.Errorf("error parsing buffer : %v", err)
			}

			readToIndex -= parsedBytes
			r.RequestLine = *reqLine
		}
	}

	//fmt.Printf("parsedBytesInt: %q\n", parsedBytesInt)
	return r, nil
}

func ParseRequestLineFromString(s string) (*RequestLine, error) {
	//fmt.Println("S", s)
	trimed := strings.TrimSuffix(s, "\r\n")
	//fmt.Println("TrimSuffix", trimed)
	splited := strings.Split(trimed, " ")
	for i, s := range splited {
		fmt.Printf("splited[%d]: %q\n", i, s)
	}
	if len(splited) != 3 {
		return nil, fmt.Errorf("Bad request line, lengt has to be 3, provided length is %d", len(splited))
	}

	method := string(splited[0])
	fmt.Printf("method: %q\n", method)

	for i, c := range method {
		fmt.Printf("index=%d char=%c ascii=%d\n", i, c, c)
	}
	for _, c := range method {
		if rune(c) < 'A' || rune(c) > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := splited[1]
	fmt.Printf("requestTarget: %q\n", requestTarget)

	httpVersion := strings.Split(splited[2], "/")
	fmt.Printf("httpVersion: %q\n", httpVersion[1])
	if len(httpVersion) != 2 {
		return nil, fmt.Errorf("Malformed starline... length is %d (has to be 2)", len(httpVersion))
	}

	if httpVersion[0] != "HTTP" {
		return nil, fmt.Errorf("Bad HTTP version, want HTTP, provided: %s ", httpVersion[0])
	}
	if string(httpVersion[1]) != "1.1" {
		return nil, fmt.Errorf("Bad HTTP version, want 1.1, provided: %s ", httpVersion[1])
	}

	return &RequestLine{
		HttpVersion:   httpVersion[1],
		RequestTarget: requestTarget,
		Method:        method,
	}, nil
}
