package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"log"
	"math/big"
	"net/http"

	"github.com/quic-go/quic-go/http3"
)

func main() {
	server := http3.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("received request")
			w.WriteHeader(http.StatusOK)
		}),
		Addr:      ":4433",
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{tlsCert()}},
	}
	server.ListenAndServe()
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
