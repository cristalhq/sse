package sse

import (
	"net/http"
	"strconv"
	"time"
)

var defaultUpgrader Upgrader

func UpgradeHTTP(r *http.Request, w http.ResponseWriter) (*Stream, error) {
	return defaultUpgrader.UpgradeHTTP(r, w)
}

var noDeadline time.Time

type Upgrader struct {
	Timeout time.Duration
}

func (u Upgrader) UpgradeHTTP(r *http.Request, w http.ResponseWriter) (*Stream, error) {
	fl, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Flushing not supported", http.StatusNotImplemented)
		return nil, ErrNotFlusher
	}

	h := w.Header()
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("Content-Type", "text/event-stream")
	h.Set("Transfer-Encoding", "chunked")

	w.WriteHeader(http.StatusOK)
	fl.Flush() // flush headers

	s := &Stream{
		w:       w,
		flusher: fl,
	}
	return s, nil
}

// httpError is like the http.Error with additional headers.
func httpError(w http.ResponseWriter, body string, code int) {
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(body))
}
