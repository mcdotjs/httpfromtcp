package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const CRLF string = "\r\n"

func ParseRequestLine(b []byte) (*RequestLine, error) {
	idxOfRegisteredNurse := bytes.Index(b, []byte(CRLF))
	if idxOfRegisteredNurse == -1 {
		return &RequestLine{}, fmt.Errorf("%s", "Cannot find first registered nurse")
	}
	firstLine := string(b[:idxOfRegisteredNurse])
	splited := strings.Split(firstLine, " ")
	if len(splited) != 3 {
		return nil, fmt.Errorf("Bad request line, lengt has to be 3, provided length is %d", len(splited))
	}

	method := splited[0]
	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return nil, fmt.Errorf("Bad request line mehod %s", method)
		}
	}

	requestTarget := splited[1]

	httpVersion := strings.Split(splited[2], "/")

	if len(httpVersion) != 2 {
		return nil, fmt.Errorf("Malformed starline... length is %d (has to be 2)", len(httpVersion))
	}

	if httpVersion[0] != "HTTP" {
		return nil, fmt.Errorf("Bad HTTP version, want HTTP, provided: %s ", httpVersion[0])
	}
	if httpVersion[1] != "1.1" {
		return nil, fmt.Errorf("Bad HTTP version, want 1.1, provided: %s ", httpVersion[1])
	}

	return &RequestLine{
		HttpVersion:   httpVersion[1],
		RequestTarget: requestTarget,
		Method:        method,
	}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{}
	b, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, err
	}

	requstLine, err := ParseRequestLine(b)
	if err != nil {
		return nil, err
	}

	r.RequestLine = *requstLine
	return r, nil
}
