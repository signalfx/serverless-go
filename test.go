// Package p contains an HTTP Cloud Function.
package sfxserverless

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	"github.com/seonsfx/serverless-go/gcfwrapper"
)

var wrapper gcfwrapper.HandlerWrapper

func handler(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Hello World!")
		return
	}
	if d.Message == "" {
		fmt.Fprint(w, "Hello World!")
		return
	}
	fmt.Fprint(w, html.EscapeString(d.Message))
}

// Test prints the JSON encoded "message" field in the body
// of the request or "Hello, World!" if there isn't one.
func Test(w http.ResponseWriter, r *http.Request) {
	wrapper = gcfwrapper.NewHandlerWrapper(handler)
	wrapper.Invoke(w, r)
}
