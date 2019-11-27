package sse

import (
	"bufio"
	"bytes"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	s, buf := newStream()

	if err := s.SetRetry(1234 * time.Second); err != nil {
		t.Fatalf("unexpected err %#v", err)
	}
	if err := s.Flush(); err != nil {
		t.Fatalf("unexpected err %#v", err)
	}

	gotBytes := buf.Bytes()
	got := string(gotBytes)

	if want := "retry:1234000\n"; got != want {
		t.Fatalf("want %#v, got %#v", want, got)
	}
}

func newStream() (*Stream, *bytes.Buffer) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	r := bufio.NewReader(bytes.NewReader(make([]byte, 0, 1024)))
	w := bufio.NewWriter(buf)
	bw := bufio.NewReadWriter(r, w)

	s := &Stream{
		bw: bw,
		w:  w,
	}
	return s, buf
}
