package sse

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
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

func TestStream_Close(t *testing.T) {
	server := newServer(func(w http.ResponseWriter, r *http.Request) {
		u := Upgrader{}

		stream, err := u.UpgradeHTTP(r, w)
		if err != nil {
			t.Fatal(err)
		}

		if err := stream.Close(); err != nil {
			t.Fatalf("stream.Close() = %v, want nil", err)
		}
	})

	client := http.DefaultClient

	resp, err := client.Do(newStreamRequest(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	want := "event:close\ndata:\n\n"

	if got := string(body); got != want {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}
