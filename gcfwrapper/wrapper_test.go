package gcfwrapper

import (
	"encoding/json"
	"net/http"
	"testing"
)

var wrapper HandlerWrapper

func handlerThrowsError(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Message string `json:"message"`
	}

	// Should throw an error because r is missing Body
	json.NewDecoder(r.Body).Decode(&d)
}

// TestErrorHandling should throw an error from the user function
// to see if it is covered by the wrapper
func TestErrorHandling(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("User function error was not handled.")
		}
	}()

	var w http.ResponseWriter
	r := &http.Request{}
	wrapper = NewHandlerWrapper(handlerThrowsError)
	wrapper.Invoke(w, r)
}
