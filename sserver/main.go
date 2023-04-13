package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"

	"github.com/blck-snwmn/hello-quicgo/schema/fbs"
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
		go handleConn(conn)
	}
}

func handleConn(conn quic.Connection) {
	fmt.Println("recieved connection")
	defer fmt.Println("close connection")
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		panic(err)
	}

	for {
		fmt.Printf("streamID=%d\n", stream.StreamID())

		var lengthBuf [4]byte
		_, err := stream.Read(lengthBuf[:])
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Printf("failed to read length: %v\n", err)
			return
		}
		length := binary.LittleEndian.Uint32(lengthBuf[:])
		fmt.Printf("length=%d\n", length)

		buf := make([]byte, length)
		_, err = stream.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Printf("failed to read message: %v\n", err)
			return
		}

		user := fbs.GetRootAsUser(buf, 0)

		var pos fbs.Position
		user.Pos(&pos)

		fmt.Printf("name=`%s`, ", user.Name())
		fmt.Printf("color=%v, ", user.Color())
		fmt.Printf("position{x, y, z} = {%v, %v, %v}, ", pos.X(), pos.Y(), pos.Z())

		for i := 0; i < user.InventoryLength(); i++ {
			var item fbs.Item
			user.Inventory(&item, i)
			fmt.Printf("item[%d]:name=%s,", i, item.Name())
		}
		fmt.Println()
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
