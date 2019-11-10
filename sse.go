package sse

import (
	"net/http"
	"strconv"
	"time"
)

var defaultUpgrader = Upgrader{
	Autoflush: true,
}

func UpgradeHTTP(r *http.Request, w http.ResponseWriter) (*Stream, error) {
	return defaultUpgrader.UpgradeHTTP(r, w)
}

var noDeadline time.Time

type Upgrader struct {
	Timeout   time.Duration
	Autoflush bool
}

func (u Upgrader) UpgradeHTTP(r *http.Request, w http.ResponseWriter) (*Stream, error) {
	fl, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Flushing not supported", http.StatusNotImplemented)
		return nil, ErrNotFlusher
	}

	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")

	w.WriteHeader(http.StatusOK)
	fl.Flush() // flush headers

	s := &Stream{
		w:         w,
		flusher:   fl,
		autoFlush: u.Autoflush,
	}
	return s, nil
}

// httpError is like the http.Error with additional headers.
func httpError(w http.ResponseWriter, body string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(code)
	w.Write([]byte(body))
}
