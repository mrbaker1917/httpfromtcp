package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func IsAllUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func parseRequestLine(r string) (RequestLine, error) {
	rSplit := strings.Split(r, " ")
	if len(rSplit) != 3 {
		return RequestLine{}, errors.New("Request input string not correct")
	}
	method := rSplit[0]
	if !IsAllUpper(method) {
		return RequestLine{}, errors.New("Request method not correct format")
	}
	rTarget := rSplit[1]
	hversionParts := strings.Split(rSplit[2], "/")
	if hversionParts[0] != "HTTP" || hversionParts[1] != "1.1" {
		return RequestLine{}, errors.New("Request is not HTTP!")
	}

	hversion := hversionParts[1]
	rLine := RequestLine{
		HttpVersion:   hversion,
		RequestTarget: rTarget,
		Method:        method,
	}
	return rLine, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	reqParts := strings.Split(string(b), "\r\n")

	reqLine, err := parseRequestLine(reqParts[0])
	if err != nil {
		return nil, err
	}
	req := Request{
		RequestLine: reqLine,
	}
	return &req, nil
}
