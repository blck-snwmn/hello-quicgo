package main

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

func main() {
	s := webtransport.Server{
		H3: http3.Server{Addr: ":4433"},
	}

	// Create a new HTTP endpoint /webtransport.
	http.HandleFunc("/webtransport", func(w http.ResponseWriter, r *http.Request) {
		log.Println("recieve")
		defer log.Println("done")
		conn, err := s.Upgrade(w, r)
		if err != nil {
			log.Printf("upgrading failed: %s", err)
			w.WriteHeader(500)
			return
		}

		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			log.Printf("opening stream failed: %s", err)
			w.WriteHeader(500)
			return
		}
		io.Copy(stream, stream)
		// Handle the connection. Here goes the application logic.
	})

	err := s.ListenAndServeTLS("./testdata/server.crt", "./testdata/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
