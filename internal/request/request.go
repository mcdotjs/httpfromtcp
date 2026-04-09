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
	Initialized = iota
	Done
)

func ParseRequestLine(b []byte) (int, error) {
	//TODO:
	//Update your parseRequestLine to return the number of bytes it consumed.
	//If it can't find an \r\n (this is important!) it should return 0 and no error.
	//This just means that it needs more data before it can parse the request line.

	idxOfRegisteredNurse := bytes.Index(b, []byte(CRLF))
	if idxOfRegisteredNurse == -1 {
		//NOTE: ak nenajde crlf .... just give me more data :)
		return 0, nil
	}

	return idxOfRegisteredNurse, nil
}

func (r *Request) parse(data []byte) (int, error) {
	//TODO:
	//It accepts all currently unparsed bytes from the buffer
	//It updates the "state" of the parser, and the parsed RequestLine field.
	//It returns the number of bytes it consumed (meaning successfully parsed) and an error if it encountered one.
	fmt.Printf("parse -> data: %q\n", data)
	if len(data) == 0 {
		return 0, nil
	}

	if r.ParserState == Initialized {
		// bytesConsumed, err := ParseRequestLine(data)
		// fmt.Println("PARSE bytesConsumed: ", bytesConsumed)
		// if err == nil && bytesConsumed == 0 {
		// 	return bytesConsumed, nil
		// }
		//
		// if err != nil {
		// 	return bytesConsumed, err
		// }

		firstLine := string(data)

		fmt.Println("RRRR!!: ", firstLine)
		reqLine, err := ParseRequestLineFromString(firstLine)
		fmt.Println("RRRR: ", reqLine)
		if err != nil {
			return 0, fmt.Errorf("Error parsing from string: %v", err)
		}
		r.RequestLine = *reqLine
		return 0, nil
	} else if r.ParserState == Done {
		return 0, fmt.Errorf("error trying parse on done state")
	} else {
		return 0, fmt.Errorf("something error parse")
	}
	//returning parsed number of bytes
}


// TODO:
// Instead of reading all the bytes, and then parsing the request line,
// it should use a loop to continually read from the reader and parse new chunks using the parse method.
// The loop should continue until the parser is in the "done" state.
func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{
		ParserState: Initialized,
	}
	bufferSize := 8
	readToIndex := 0
	parsedBytes1 := 0
	buffer := make([]byte, bufferSize)
	for {
		if r.ParserState != Done {

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
			fmt.Printf("buffer3: %q\n", buffer)
			if err != nil {
				if err == io.EOF {
					r.ParserState = Done
					break
				}
				return nil, fmt.Errorf("error reading to buffer")
			}

			fmt.Printf("buffer4: %q\n", buffer)
			readToIndex += readBytes

			parsedBytes, err := ParseRequestLine(buffer)

			fmt.Printf("if 0 need more data: %d\n", parsedBytes)
			if parsedBytes != 0 {
				parsedBytes1 = parsedBytes
				break
			}
		}
	}

	fmt.Printf("vonku")
	//parsedBytesInt, err := r.parse(buffer[:parsedBytes1])

	reqLine, err := ParseRequestLineFromString(string(buffer[:parsedBytes1]))
	fmt.Println("RE:", reqLine)
	if err != nil {
		return nil, fmt.Errorf("error parsing buffer : %v", err)
	}

	r.RequestLine = *reqLine
	//fmt.Printf("parsedBytesInt: %q\n", parsedBytesInt)
	return r, nil
}

func RequestFromReader1(reader io.Reader) (*Request, error) {
	//TODO:
	//Instead of reading all the bytes, and then parsing the request line,
	//it should use a loop to continually read from the reader and parse new chunks using the parse method.
	//The loop should continue until the parser is in the "done" state.
	r := &Request{
		ParserState: Initialized,
	}
	bufferSize := 8
	readToIndex := 0
	buffer := make([]byte, bufferSize)
	//NOTE: nie citat all data
	//ale loopovat s pouzitim parse method until done state
	for {
		if r.ParserState != Done {
			if len(buffer) < readToIndex {
				fmt.Println("growing", len(buffer), readToIndex)
				bufferSize *= 2
				temp := make([]byte, bufferSize)
				copy(temp, buffer)
				buffer = temp
			}
			tempBuffer := make([]byte, bufferSize)
			readBytes, err := reader.Read(tempBuffer)

			buffer = append(buffer, tempBuffer...)
			readToIndex += readBytes

			//NOTE: remove null bytes from beggining
			parsedBytesInt, err := r.parse(buffer[readBytes:])
			fmt.Println("cutt:::", parsedBytesInt)
			if err == nil && parsedBytesInt == 0 {
				fmt.Println("::::::::::::::::::::::::: ", len(buffer), string(buffer))
				fmt.Printf("more data: %q\n", tempBuffer)
				fmt.Printf("more date readToIndex: %d\n", readToIndex)
				continue
			}
			if err == nil && parsedBytesInt != 0 {
				fmt.Printf("parsedBytesInt neni0: %d\n", parsedBytesInt)
			}
			if err != nil {
				if err == io.EOF {
					bufferSize = 0
					buffer = make([]byte, bufferSize)
					r.ParserState = Done
					break
				}
			}

			readToIndex -= parsedBytesInt
			return r, err
		}
		// b, err := r.parse(buffer[readToIndex:])
		// if err != nil {
		// 	return r, nil
		// }
		// readToIndex -= b
	}
	return r, nil
}

func ParseRequestLineFromString(s string) (*RequestLine, error) {
	fmt.Println("S", s)
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
	// for _, char := range method {
	// 	if char < 'A' || char > 'Z' {
	// 		fmt.Println("char: ", char)
	// 		return nil, fmt.Errorf("Bad request line mehod %s", method)
	// 	}
	// }

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
