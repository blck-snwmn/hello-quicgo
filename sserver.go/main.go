package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"
	"log"
	"math/big"

	"github.com/quic-go/quic-go"
)

func main() {
	err := s()
	if err != nil {
		log.Fatal(err)
	}
}

func s() error {
	addr := ":4433"
	c := &tls.Config{
		Certificates: []tls.Certificate{tlsCert()},
		NextProtos:   []string{"quic-echo-example"},
	}
	listener, err := quic.ListenAddr(addr, c, nil)
	if err != nil {
		return err
	}
	conn, err := listener.Accept(context.Background())
	if err != nil {
		return err
	}
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		panic(err)
	}
	// Echo through the loggingWriter
	_, err = io.Copy(stream, stream)
	return err
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
