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
	n = len(headerLineText) + 2
	key, val, found := strings.Cut(headerLineText, ":")
	if !found {
		return 0, false, errors.New("data is not in key-value format")
	}
	keyTrim := strings.TrimRight(key, " ")
	if len(keyTrim) != len(key) {
		return 0, false, errors.New("No spaces allowed in the field name/key")
	}
	if !isValidHeaderKey(keyTrim) {
		return 0, false, errors.New("invalid characters in header key")
	}

	val = strings.TrimSpace(val)
	key = strings.TrimSpace(key)
	key = strings.ToLower(key)

	if _, ok := h[key]; ok {
		preval := h[key]
		h[key] = preval + ", " + val
		return n, false, nil
	} else {
		h[key] = val
	}

	return n, false, nil
}

func isValidHeaderKey(key string) bool {
	allowed := " ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!#$%&'*+-.^_`|~"
	for _, c := range key {
		if !strings.ContainsRune(allowed, c) {
			return false
		}
	}
	return true
}
