package apps

import (
	"fmt"
	"io"
	"net/http"

	"github.com/yunabe/gae-codelab/mylib"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		name = "World"
	}
	io.WriteString(w, mylib.GetHelloMessage(name))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("No handler is registere for %s", r.URL.Path), http.StatusOK)
}

func init() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/", defaultHandler)
}
