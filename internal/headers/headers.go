package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

const CRLF string = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	stringData := string(data)
	fmt.Printf("data: %q \n", stringData)

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
	colonIdx := bytes.Index(data, []byte(":"))
	if colonIdx == -1 {
		return 0, false, fmt.Errorf("nemame colon index")
	}

	// parts := bytes.SplitN(data[:crlfIdx], []byte(":"), 2)
	//
	// fmt.Println("PARTS:", string(parts[0]), string(parts[1]))

	key := string(data[:colonIdx])

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("Key has space")
	}
	value := string(data[colonIdx+1 : crlfIdx])
	h[key] = strings.TrimSpace(value)

	// l ... fmt.Println("LLLLLLL:", len(value)+len(key), crlfIdx)
	//return l + 1 + 2, false, nil
	return crlfIdx + 2, false, nil
}
