package main

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"github.com/blck-snwmn/hello-quicgo/schema/fbs"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

func main() {
	var sg sync.WaitGroup
	for _, name := range []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"} {
		sg.Add(1)
		go func(name string) {
			defer sg.Done()
			err := exec(name)
			if err != nil {
				panic(err)
			}
		}(name)
	}
	sg.Wait()
}

func exec(name string) error {
	d := webtransport.Dialer{
		RoundTripper: &http3.RoundTripper{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	_, conn, err := d.Dial(context.Background(), "https://localhost:4433/webtransport?room=test", nil)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return fmt.Errorf("failed to open stream: %w", err)
	}
	defer stream.Close()

	tick := time.NewTicker(time.Second)
	after := time.After(10 * time.Second)

end:
	for {
		select {
		case <-after:
			break end
		case <-tick.C:
			{
				builder := flatbuffers.NewBuilder(200)
				u := fbs.UserT{
					Name: name,
					Pos: &fbs.PositionT{
						X: float32(11),
						Y: float32(12),
						Z: float32(13),
					},
				}
				builder.FinishSizePrefixed(u.Pack(builder))
				buf := builder.FinishedBytes()

				_, err = stream.Write(buf)
				if err != nil {
					return fmt.Errorf("failed to write: %w", err)
				}
			}
			{
				var lengthBuf [4]byte
				_, err = stream.Read(lengthBuf[:])
				if err != nil {
					return fmt.Errorf("failed to read length: %w", err)
				}
				length := binary.LittleEndian.Uint32(lengthBuf[:])
				fmt.Printf("length=%d\n", length)

				buf := make([]byte, length)
				_, err = stream.Read(buf)
				if err != nil {
					return fmt.Errorf("failed to read message: %w", err)
				}
				bs := fbs.GetRootAsBroadcast(buf, 0)
				bb := bs.UnPack()
				for _, u := range bb.Poss {
					fmt.Printf("[%s]{x,y,z}={%f,%f,%f}", u.Name, u.Pos.X, u.Pos.Y, u.Pos.Z)
				}
				fmt.Println()
			}
		}
	}
	return nil
}
