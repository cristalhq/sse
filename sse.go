package sse

import (
	"net/http"
	"time"
)

var defaultUpgrader Upgrader

func UpgradeHTTP(r *http.Request, w http.ResponseWriter) (*Stream, error) {
	return defaultUpgrader.UpgradeHTTP(r, w)
}

type Upgrader struct {
	Timeout time.Duration
}

func (u Upgrader) UpgradeHTTP(r *http.Request, w http.ResponseWriter) (*Stream, error) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Cannot hijack a connection", http.StatusBadRequest)
		return nil, ErrNotHijacker
	}

	_, bw, err := hj.Hijack()
	if err != nil {
		http.Error(w, http.ErrHijacked.Error(), http.StatusInternalServerError)
		return nil, http.ErrHijacked
	}

	httpWriteResponseUpgrade(bw.Writer)
	if err := bw.Flush(); err != nil {
		return nil, err
	}

	s := &Stream{
		bw: bw,
		w:  w,
	}
	return s, nil
}
