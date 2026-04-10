package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}
	headerLineText := string(data[:idx])
	key, val, found := strings.Cut(headerLineText, ":")
	if !found {
		return 0, false, errors.New("data is not in key-value format")
	}
	keyTrim := strings.TrimSpace(key)
	if len(keyTrim) != len(key) {
		return 0, false, errors.New("No spaces allowed in the field name/key")
	}
	val = strings.TrimSpace(val)
	h[key] = val
	return idx + 2, false, nil
}
