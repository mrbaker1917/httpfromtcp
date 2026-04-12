package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mrbaker1917/httpfromtcp/internal/headers"
)

type requestState int

const (
	initialized requestState = iota
	requestStateParsingHeaders
	done
)

type Request struct {
	RequestLine RequestLine
	state       requestState
	Headers     headers.Headers
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	req := &Request{
		state:   initialized,
		Headers: headers.NewHeaders(),
	}
	for req.state != done {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		// read into the buffer
		n, err := reader.Read(buf[readToIndex:])

		if errors.Is(err, io.EOF) {
			if req.state != done {
				return nil, errors.New("reached end of file, but header not terminated.")
			}
			break
		}
		if err != nil {
			return nil, err
		}
		readToIndex += n

		// parse from the buffer
		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != done {

		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, fmt.Errorf("Unable to parse line: %s", err)
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil

}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case initialized:
		requestLine, num, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if num == 0 {
			return 0, nil
		} else {
			r.RequestLine = *requestLine
			r.state = requestStateParsingHeaders
			return num, nil
		}
	case requestStateParsingHeaders:
		n, d, err := r.Headers.Parse(data)
		if err != nil {
			return 0, errors.New("error: trying to read headers")
		}

		if d {
			r.state = done
			return n, nil
		}
		return n, nil
	case done:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}
}
