package sse

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {
	const msg = "test-ok\n"

	server := newServer(func(w http.ResponseWriter, r *http.Request) {
		u := Upgrader{Autoflush: true}

		stream, err := u.UpgradeHTTP(r, w)
		if err != nil {
			t.Fatal(err)
		}

		err = stream.WriteRaw([]byte(msg))
		if err != nil {
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

	bs, err := br.ReadBytes('\n')
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != msg {
		t.Fatalf("got %#v, want %#v", string(bs), msg)
	}
}

func newServer(h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(h))
}

func newStreamRequest(url string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "text/event-stream")
	return req
}
