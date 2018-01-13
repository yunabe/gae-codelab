package apps

import (
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	mrand "math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/yunabe/gae-codelab/datastore"
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

func helloSlow(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	secStr := r.FormValue("sec")
	sec, err := strconv.ParseInt(secStr, 0, 64)
	if err != nil {
		log.Infof(ctx, "Failed to parse %q: %v", secStr, err)
		sec = 5
	}
	time.Sleep(time.Duration(sec) * time.Second)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "Waited for %d secs", sec)
}

type mathRandReader struct{}

func (mathRandReader) Read(p []byte) (int, error) {
	return mrand.Read(p)
}

func mathRandHandler(w http.ResponseWriter, r *http.Request) {
	// math/rand.Read is deterministic. This handler returns the same value
	// (52fdfc072182654f163f5f0f9a621d729566c74d10037c4d7bbb0407d1e2c649) after reload.
	buf := make([]byte, 32)
	_, err := io.ReadFull(mathRandReader{}, buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	io.WriteString(w, hex.EncodeToString(buf))
}

func cryptoRandHandler(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, 32)
	_, err := io.ReadFull(crand.Reader, buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	io.WriteString(w, hex.EncodeToString(buf))
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
	datastore.RegisterHandlers()
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/mrand", mathRandHandler)
	http.HandleFunc("/crand", cryptoRandHandler)
	http.HandleFunc("/slow", helloSlow)
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/register_sample_task", registerSampleTask)
	http.HandleFunc("/admin/sample_task", handleSampleTask)
	http.HandleFunc("/admin/cron_task", handleCronTask)
}
