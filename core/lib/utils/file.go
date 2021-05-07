package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/macinnir/dvc/core/lib/utils/errors"
)

const (
	chunkSize = 1024
)

// OpenFile opens a file in chunks
func OpenFile(filePath string) (byteCount int, buffer *bytes.Buffer, e error) {

	var part []byte
	var data *os.File
	var count int

	data, e = os.Open(filePath)
	if e != nil {
		return
	}
	defer data.Close()

	reader := bufio.NewReader(data)
	buffer = bytes.NewBuffer(make([]byte, 0))
	part = make([]byte, chunkSize)

	for {
		if count, e = reader.Read(part); e != nil {
			break
		}
		buffer.Write(part[:count])
	}
	if e != io.EOF {
		e = errors.NewError(fmt.Sprint("Error Reading ", filePath, ": ", e))
		return
	}

	e = nil

	byteCount = buffer.Len()
	return
}
