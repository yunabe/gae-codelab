package apps

import (
	"io/ioutil"
	"net/http"
)

import "google.golang.org/appengine/log"
import "google.golang.org/appengine"
import "google.golang.org/appengine/taskqueue"

func registerSampleTask(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	_, err := taskqueue.Add(ctx, &taskqueue.Task{
		Path:    "/admin/sample_task",
		Payload: []byte("Hello taskqueue!"),
	}, "sample-task")
	if err != nil {
		log.Errorf(ctx, "Failed to register a task: %v", err)
	}
}

func handleSampleTask(w http.ResponseWriter, r *http.Request) {
	// TODO: Protect the method from CSRF.
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "handleSampleTask is invoked.")
	if r.Method != http.MethodPost {
		log.Warningf(ctx, "Received a request with unexpected method: %v", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(ctx, "Failed to load payload: %v", err)
		return
	}
	log.Infof(ctx, "Payload in string: %v", string(payload))
	for key, values := range r.Header {
		// Task API sets headers with "X-Appengine-Task" prefix.
		log.Infof(ctx, "Header[%q] = %v", key, values)
	}

	w.WriteHeader(http.StatusInternalServerError)
}
