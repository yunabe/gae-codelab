package spa

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Cache-Control", "max-age=0")
	io.WriteString(w, fmt.Sprintf("Hello SPA (%s)", time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006")))
}

func init() {
	http.HandleFunc("/api/hello", helloHandler)
}
