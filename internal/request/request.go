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

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	crlfIdx := bytes.Index(b, []byte(CRLF))
	if crlfIdx == -1 {
		//NOTE: ak nenajde crlf .... just give me more data :)
		return nil, 0, nil
	}
	requestLineText := string(b[:crlfIdx])
	fmt.Println("requst line text: ", requestLineText)
	reqLine, err := parseRequestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return reqLine, crlfIdx + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.ParserState {
	case initialized:
		reqLine, parsedBytes, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if parsedBytes == 0 {
			return 0, nil
		}
		r.RequestLine = *reqLine
		r.ParserState = done
		return parsedBytes, nil
	case done:
		return 0, fmt.Errorf("error: parsing request is done")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{
		ParserState: initialized,
	}
	bufferSize := 8
	readToIndex := 0
	buffer := make([]byte, bufferSize)
	for r.ParserState != done {

		if len(buffer) == readToIndex {
			//fmt.Printf("buffer growing: %q\n", buffer)
			bufferSize *= 2
			temp := make([]byte, bufferSize)
			copy(temp, buffer)
			buffer = temp
			//fmt.Printf("buffer growing1: %q\n", buffer)
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

		//	fmt.Printf("buffer: %q\n", buffer)
		readToIndex += readBytes
		parsedBytes, err := r.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		fmt.Printf("buffer2: %q\n", buffer)
		fmt.Printf("buffer2!!!: %q\n", buffer[parsedBytes:])
		copy(buffer, buffer[parsedBytes:])
		//NOTE: check notes
		//buffer = buffer[parsedBytes:]
		fmt.Printf("buffer3: %q\n", buffer)
		readToIndex -= parsedBytes
	}

	//fmt.Printf("parsedBytesInt: %q\n", parsedBytesInt)
	fmt.Println("___________________________________________________________________________")
	return r, nil
}

func parseRequestLineFromString(s string) (*RequestLine, error) {
	parts := strings.Split(s, " ")
	for i, s := range parts {
		fmt.Printf("parts[%d]: %q\n", i, s)
	}
	if len(parts) != 3 {
		return nil, fmt.Errorf("Bad request line, lengt has to be 3, provided length is %d", len(parts))
	}

	method := string(parts[0])
	//fmt.Printf("method: %q\n", method)

	// for i, c := range method {
	// 	fmt.Printf("index=%d char=%c ascii=%d\n", i, c, c)
	// }
	for _, c := range method {
		if rune(c) < 'A' || rune(c) > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]
	fmt.Printf("requestTarget: %q\n", requestTarget)

	httpVersion := strings.Split(parts[2], "/")
	//fmt.Printf("httpVersion: %q\n", httpVersion[1])
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
