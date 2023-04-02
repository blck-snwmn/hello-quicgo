package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net/http"

	"github.com/quic-go/quic-go/http3"
)

func main() {
	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: true,
			KeyLogWriter:       io.Discard,
		},
	}
	hclient := &http.Client{
		Transport: roundTripper,
	}
	hclient.Get("https://localhost:4433")
}
