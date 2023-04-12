package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/blck-snwmn/hello-quicgo/schema/fbs"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/quic-go/quic-go"
)

var (
	addr    = "localhost:4433"
	tlsConf = &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
)

func main() {
	err := startClient()
	if err != nil {
		log.Printf("faild to startClient: %v", err)
	}
}

func startClient() error {
	conn, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}
	defer stream.Close()

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)

		builder := flatbuffers.NewBuilder(200)
		u := fbs.UserT{
			Name: "John Doe",
			Pos: &fbs.PositionT{
				X: float32(11 + i),
				Y: float32(12 + i),
				Z: float32(13 + i),
			},
			Color: fbs.Color(10),
			Inventory: []*fbs.ItemT{
				{Name: "sword"},
				{Name: "shield"},
				{Name: "armor"},
			},
		}

		builder.FinishSizePrefixed(u.Pack(builder))
		buf := builder.FinishedBytes()

		fmt.Printf("[%d]length=%d, msg=`%X`\n", i, len(buf), buf)

		_, err := stream.Write(buf)
		if err != nil {
			fmt.Printf("faild to write message: %v", err)
		}
	}

	return nil
}
