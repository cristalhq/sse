package sse

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUpgrader(t *testing.T) {
	const msg = "test-ok\n"

	server := newServer(func(w http.ResponseWriter, r *http.Request) {
		u := Upgrader{}

		stream, err := u.UpgradeHTTP(r, w)
		if err != nil {
			t.Fatal(err)
		}

		if err := stream.SetRetry(1234 * time.Second); err != nil {
			t.Fatal(err)
		}
		if err := stream.SetID(6789); err != nil {
			t.Fatal(err)
		}
		if err := stream.SetEvent("test-event"); err != nil {
			t.Fatal(err)
		}
		if err := stream.WriteBytes([]byte(msg)); err != nil {
			t.Fatal(err)
		}
		if err := stream.Flush(); err != nil {
			t.Fatal(err)
		}
	})

	client := http.DefaultClient

	resp, err := client.Do(newStreamRequest(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	br := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	wantLines := []string{
		"retry:1234000\n",
		"id:6789\n",
		"event:test-event\n",
		"data:test-ok\n",
	}
	for _, want := range wantLines {
		bs, err := br.ReadBytes('\n')
		if err != nil {
			t.Fatal(err)
		}
		if string(bs) != want {
			t.Fatalf("got %#v, want %#v", string(bs), want)
		}
	}
}

func newServer(h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(h))
}

func newStreamRequest(url string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	return req
}
