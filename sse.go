package sse

import (
	"net/http"
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

	hj, ok := w.(http.Hijacker)
	if !ok {
		httpError(w, ErrNotHijacker.Error(), http.StatusInternalServerError)
		return nil, ErrNotHijacker
	}

	conn, _, err := hj.Hijack()
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	// Clear deadlines set by server.
	conn.SetDeadline(noDeadline)
	if u.Timeout != 0 {
		conn.SetWriteDeadline(time.Now().Add(u.Timeout))
		defer conn.SetWriteDeadline(noDeadline)
	}

	h := w.Header()
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("Content-Type", "text/event-stream")

	fl.Flush() // flush headers

	s := &Stream{
		conn:      conn,
		w:         w,
		flusher:   fl,
		autoFlush: u.Autoflush,
	}
	return s, nil
}
