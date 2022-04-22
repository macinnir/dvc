package forms

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/macinnir/dvc/core/lib/utils/request"
)

func FormToBytes(r *request.Request) ([]byte, error) {
	var e error

	// Max upload of 10 mb
	if e := r.RootRequest.ParseMultipartForm(10 << 20); e != nil {
		return nil, fmt.Errorf("ParseMultiPartForm: %w", e)
	}

	var file multipart.File
	var header *multipart.FileHeader

	file, header, e = r.RootRequest.FormFile("myFile")

	if e != nil {
		return nil, e
	}

	// Create a buffer to store the header of the file in
	fileHeader := make([]byte, 512)

	// Copy the headers into the FileHeader buffer
	if _, e := file.Read(fileHeader); e != nil {
		return nil, fmt.Errorf("read FileHeader: %w", e)
	}

	// Set Position back to start
	if _, e := file.Seek(0, 0); e != nil {
		return nil, fmt.Errorf("upload.Seek(): %w", e)
	}

	fmt.Printf("Name: %#v\n", header.Filename)
	// fmt.Printf("Size: %#v\n", upload.(Sizer).Size())
	fmt.Printf("MIME: %#v\n", http.DetectContentType(fileHeader))

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	buf := bytes.NewBuffer(nil)
	if _, e = io.Copy(buf, file); e != nil {
		return nil, e
	}

	return buf.Bytes(), e
}
