package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"

	"github.com/quic-go/quic-go"
)

var (
	addr    = ":4433"
	tlsConf = &tls.Config{
		Certificates: []tls.Certificate{tlsCert()},
		NextProtos:   []string{"quic-echo-example"},
	}
)

func main() {
	err := startServer()
	if err != nil {
		log.Fatal(err)
	}
}

func startServer() error {
	addr := ":4433"
	listener, err := quic.ListenAddr(addr, tlsConf, nil)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			return fmt.Errorf("failed to accept: %v", err)
		}
		go func() {
			fmt.Println("recieved connection")
			stream, err := conn.AcceptStream(context.Background())
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(loggingWriter{stream}, stream)
			if err != nil {
				log.Printf("failed to copy: %v\n", err)
				return
			}
		}()
	}
}

type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
}

func tlsCert() tls.Certificate {
	key, _ := rsa.GenerateKey(rand.Reader, 1024)
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, _ := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	tlsCert, _ := tls.X509KeyPair(certPEM, keyPEM)
	return tlsCert
}
