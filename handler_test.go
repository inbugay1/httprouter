package httprouter_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inbugay1/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestHandlerFunc_Handle(t *testing.T) {
	t.Parallel()

	handler := httprouter.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) error {
		return errors.New("some error")
	})

	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	responseWriterRecorder := httptest.NewRecorder()

	err := handler.Handle(responseWriterRecorder, request)
	assert.EqualError(t, err, "some error")
}
