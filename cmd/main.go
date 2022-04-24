package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/cristalhq/sse"
)

func f() {
	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		stream, err := sse.UpgradeHTTP(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := struct {
			Text string `json:"text"`
		}{
			Text: "hey there",
		}
		stream.WriteJSON(data)
	})
}

const timeout = 30 * time.Second

func do() {
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	var u sse.Upgrader
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		stream, err := u.Upgrade(conn)
		if err != nil {
			log.Fatal(err)
		}

		start := time.Now()
		for time.Now().Sub(start) < 2*timeout {
			var data = struct {
				Time time.Time `json:"time"`
			}{
				Time: time.Now(),
			}
			stream.WriteJSON(data)
			// stream.WriteBytes
			// io.WriteString(w, "data: ")
			// enc.Encode(&data)
			// io.WriteString(w, "\n\n")
			stream.Flush()
			time.Sleep(time.Second)
		}
		// io.WriteString(w, "event: close\ndata:\n\n")
		stream.Close()

		// err = stream.WriteRaw([]byte(msg))
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}
}

func main() {
	// do()

	a := "/sse"
	b := "/sse2"
	a, b = b, a

	http.HandleFunc(a, func(w http.ResponseWriter, r *http.Request) {
		stream, err := sse.Upgrader{}.UpgradeHTTP(r, w)
		if err != nil {
			http.Error(w, err.Error(), 503)
			return
		}

		// var data = struct {
		// 	Time time.Time `json:"time"`
		// }{
		// 	Time: time.Now(),
		// }
		// stream.WriteJSON("123", "message", data)

		// for {
		// 	stream.Flush()
		// 	time.Sleep(1234 * time.Millisecond)
		// 	println("pong")
		// }

		// enc := json.NewEncoder(w)
		start := time.Now()
		for time.Now().Sub(start) < 2*timeout {
			var data = struct {
				Time time.Time `json:"time"`
			}{
				Time: time.Now(),
			}
			stream.WriteJSON(data)
			// stream.WriteBytes
			// io.WriteString(w, "data: ")
			// enc.Encode(&data)
			// io.WriteString(w, "\n\n")
			stream.Flush()
			time.Sleep(time.Second)
		}
		// io.WriteString(w, "event: close\ndata:\n\n")
		stream.Close()
		// stream.Flush()
		log.Printf("%q: served events for %v", r.URL.Path, time.Now().Sub(start))
	})

	http.HandleFunc(b, func(w http.ResponseWriter, r *http.Request) {
		flush := func() { log.Print("warning: flush not implemented") }
		if f, ok := w.(http.Flusher); ok {
			flush = f.Flush
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		enc := json.NewEncoder(w)
		start := time.Now()
		for time.Now().Sub(start) < timeout {
			var data = struct {
				Time time.Time `json:"time"`
			}{
				Time: time.Now(),
			}
			io.WriteString(w, "data: ")
			enc.Encode(&data)
			io.WriteString(w, "\n\n")
			flush()
			time.Sleep(time.Second)
		}
		io.WriteString(w, "event: close\ndata:\n\n")
		flush()
		log.Printf("%q: served events for %v", r.URL.Path, time.Now().Sub(start))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
		log.Printf("%q: served index.html", r.URL.Path)
	})

	http.ListenAndServe(":8080", nil)
}
