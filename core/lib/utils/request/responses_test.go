package request

import (
	"errors"
	"net/http/httptest"
	"testing"
)

func TestInternalServerError(t *testing.T) {
	r := &Request{}
	w := httptest.NewRecorder()
	InternalServerError(r, w, errors.New("Test error"))
}
