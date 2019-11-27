package sse

import (
	"bufio"
	"fmt"
	"net"
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

func TestUpgrader2(t *testing.T) {
	const msg = "test-ok\n"

	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		t.Fatal(err)
	}

	var u Upgrader
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				t.Fatal(err)
			}

			stream, err := u.Upgrade(conn)
			if err != nil {
				t.Fatal(err)
			}

			err = stream.WriteBytes([]byte(msg))
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	client := http.DefaultClient

	resp, err := client.Do(newStreamRequest("http://localhost:8080"))
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

// func TestDebugUpgrader(t *testing.T) {
// 	for _, test := range []struct {
// 		upgrader Upgrader
// 		req      []byte
// 	}{
// 		{
// 			// Base case.
// 		},
// 		{
// 			req: []byte("" +
// 				"GET /test HTTP/1.1\r\n" +
// 				"Host: example.org\r\n" +
// 				"\r\n",
// 			),
// 		},
// 		{
// 			req: []byte("PUT /fail HTTP/1.1\r\n\r\n"),
// 		},
// 		{
// 			req: []byte("GET /fail HTTP/1.0\r\n\r\n"),
// 		},
// 	} {
// 		t.Run("test", func(t *testing.T) {
// 			var (
// 				reqBuf bytes.Buffer
// 				resBuf bytes.Buffer

// 				expReq, expRes []byte
// 				actReq, actRes []byte
// 			)
// 			if test.req == nil {
// 				var dialer Upgrader
// 				dialer.Upgrade(struct {
// 					io.Reader
// 					io.Writer
// 				}{
// 					new(falseReader),
// 					&reqBuf,
// 				})
// 			} else {
// 				reqBuf.Write(test.req)
// 			}

// 			// Need to save bytes before they will be read by Upgrade().
// 			expReq = reqBuf.Bytes()

// 			du := DebugUpgrader{
// 				Upgrader:   test.upgrader,
// 				OnRequest:  func(p []byte) { actReq = p },
// 				OnResponse: func(p []byte) { actRes = p },
// 			}
// 			du.Upgrade(struct {
// 				io.Reader
// 				io.Writer
// 			}{
// 				&reqBuf,
// 				&resBuf,
// 			})

// 			expRes = resBuf.Bytes()

// 			if !bytes.Equal(actReq, expReq) {
// 				t.Errorf(
// 					"unexpected request bytes:\nact:\n%s\nwant:\n%s\n",
// 					actReq, expReq,
// 				)
// 			}
// 			if !bytes.Equal(actRes, expRes) {
// 				t.Errorf(
// 					"unexpected response bytes:\nact:\n%s\nwant:\n%s\n",
// 					actRes, expRes,
// 				)
// 			}
// 		})
// 	}
// }

type falseReader struct{}

func (f falseReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("falsy read")
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
