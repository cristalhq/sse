package sse_test

import (
	"log"
	"net"
	"net/http"

	"github.com/cristalhq/sse"
)

func ExampleUpgradeHTTP() {
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

func ExampleUpgrader() {
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	var u sse.Upgrader
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept err: %#v", err)
			continue
		}

		stream, err := u.Upgrade(conn)
		if err != nil {
			log.Printf("upgrade err: %#v", err)
			continue
		}

		go func() {
			defer stream.Close()

			err := stream.WriteString(`hey there`)
			if err != nil {
				log.Printf("send err: %#v", err)
			}
		}()
	}
}
