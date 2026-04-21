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

type WriterState int

const (
	writingStatusLine WriterState = iota
	writingHeaders
	writingBody
	done
)

type Writer struct {
	writerState WriterState
	writer      io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writingStatusLine {
		return fmt.Errorf("Writing status line should be first.")
	}
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
	_, err := w.writer.Write([]byte(statusLine))
	w.writerState = writingHeaders

	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = fmt.Sprintf("%d", contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"
	return h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writingHeaders {
		return fmt.Errorf("Writing headers should be second.")
	}
	for k, v := range headers {
		headerLine := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.writer.Write([]byte(headerLine))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	w.writerState = writingBody
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writingBody {
		return 0, fmt.Errorf("Write statusLine and headers first")
	}

	n, err := w.writer.Write(p)
	if err != nil {
		return 0, fmt.Errorf("Error in writing body: %v.", err)
	}
	w.writerState = done
	return n, nil
}
