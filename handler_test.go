package httprouter_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerFunc_Handle(t *testing.T) {
	t.Parallel()

	handlerFunc := func(_ http.ResponseWriter, _ *http.Request) error {
		return errors.New("handler error")
	}

	expectedError := errors.New("handler error")

	err := handlerFunc(nil, nil)
	assert.Equal(t, expectedError, err, "Expected error mismatch")
}
