package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

const CRLF string = "\r\n"

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func isValidString(data []byte) bool {
	for _, c := range data {
		if !isTokenChar(c) {
			return false
		}
	}
	return true
}

func isTokenChar(c byte) bool {
	if (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') {
		return true
	}
	return slices.Contains(tokenChars, c)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIdx := bytes.Index(data, []byte(CRLF))
	if crlfIdx == -1 {
		//not this return 0, false, fmt.Errorf("nemame CRLF index")
		// need more data
		return 0, false, nil
	}
	if crlfIdx == 0 {
		//If you do find a CRLF, but it's at the start of the data, you've found the end of the headers, so return the proper values immediately.
		// 2 just for /r/n
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:crlfIdx], []byte(":"), 2)

	keyString := string(parts[0])
	if keyString != strings.TrimRight(keyString, " ") {
		return 0, false, fmt.Errorf("header key has space on right: %s", keyString)
	}
	valueString := string(parts[1])

	//you have to trim key string then coerce it to []byte
	keyString = strings.TrimSpace(keyString)
	if !isValidString([]byte(keyString)) {
		return 0, false, fmt.Errorf("Key has invalid symbol")
	}

	lowerKey := strings.ToLower(keyString)
	h[lowerKey] = strings.TrimSpace(valueString)

	// l ... fmt.Println("LLLLLLL:", len(value)+len(key), crlfIdx)
	//return l + 1 + 2, false, nil
	return crlfIdx + 2, false, nil
}
