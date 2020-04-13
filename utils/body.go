package utils

import (
	"bytes"
	"io"
)

func ReadAll(r io.Reader) ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 65536))
	io.Copy(buffer, r)
	temp := buffer.Bytes()
	length := len(temp)
	var body []byte
	//are we wasting more than 10% space?
	if cap(temp) > (length + length/10) {
		body = make([]byte, length)
		copy(body, temp)
	} else {
		body = temp
	}
	return body, nil
}
