package sse

import (
	"bufio"
	"bytes"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	s, buf := newStream()

	err := s.SetRetry(1234 * time.Second)
	if err != nil {
		t.Fatalf("unexpected err %#v", err)
	}
	gotBytes := buf.Bytes()

	if want := `retry: 1234000\n`; string(gotBytes) != want {
		t.Fatalf("want %#v, got %#v", want, string(gotBytes))
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
