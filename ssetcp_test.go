package sse

import (
	"bufio"
	"net"
	"net/http"
	"testing"
	"time"
)

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
		}
	}()

	client := http.DefaultClient

	resp, err := client.Do(newStreamRequest("http://localhost:8080"))
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
