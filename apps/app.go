package apps

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/yunabe/gae-codelab/mylib"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		name = "World"
	}
	// Set content-type explicitly though ResponseWriter.Write detects it automatically.
	// https://golang.org/pkg/net/http/#ResponseWriter
	w.Header().Add("Content-Type", "text/plain")
	io.WriteString(w, mylib.GetHelloMessage(name))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("No handler is registere for %s", r.URL.Path), http.StatusOK)
}

func handleCronTask(w http.ResponseWriter, r *http.Request) {
	// TODO: Protect the method from CSRF.
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "handleCronTask is invoked with method = %q", r.Method)
	if r.Method != http.MethodGet {
		log.Warningf(ctx, "Received a request with unexpected method: %v", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(ctx, "Failed to load payload: %v", err)
		return
	}
	if len(payload) != 0 {
		log.Errorf(ctx, "Payload was not empty unexpectedly: %v", string(payload))
	}
	for key, values := range r.Header {
		// Task API sets headers with "X-Appengine-Task" prefix.
		log.Infof(ctx, "Header[%q] = %v", key, values)
	}
	w.WriteHeader(http.StatusInternalServerError)
}

func init() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/register_sample_task", registerSampleTask)
	http.HandleFunc("/admin/sample_task", handleSampleTask)
	http.HandleFunc("/admin/cron_task", handleCronTask)
}
