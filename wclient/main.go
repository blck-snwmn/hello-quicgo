package main

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

func main() {
	ctx := context.Background()
	d := webtransport.Dialer{
		RoundTripper: &http3.RoundTripper{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	_, conn, err := d.Dial(ctx, "https://localhost:4433/webtransport", nil)
	if err != nil {
		panic(err)
	}

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	const msg = "Hello, World!"
	_, err = stream.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
	b := make([]byte, len(msg))
	_, err = stream.Read(b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received: %s\n", b)
}
