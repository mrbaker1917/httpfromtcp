package response

import (
	"fmt"
	"io"

	"github.com/mrbaker1917/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	OK          StatusCode = 200
	BadRequest  StatusCode = 400
	ServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := ""
	switch statusCode {
	case 200:
		statusLine = "HTTP/1.1 200 OK\r\n"
	case 400:
		statusLine = "HTTP/1.1 400 Bad Request\r\n"
	case 500:
		statusLine = "HTTP/1.1 500 Internal Server Error\r\n"
	default:
		statusLine = fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)
	}
	_, err := w.Write([]byte(statusLine))

	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = fmt.Sprintf("%d", contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		headerLine := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.Write([]byte(headerLine))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
