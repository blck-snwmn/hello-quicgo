package main

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"sync"
	"time"

	"github.com/blck-snwmn/hello-quicgo/schema/fbs"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

func main() {
	var sg sync.WaitGroup
	for i := 0; i < 100; i++ {
		sg.Add(1)
		go func(name string) {
			defer sg.Done()
			err := exec("name-" + name)
			if err != nil {
				panic(err)
			}
		}(hex.EncodeToString([]byte{byte(i)}))
	}
	sg.Wait()
}

var (
	max = big.NewInt(100)
)

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

	tick := time.NewTicker(100 * time.Microsecond)
	after := time.After(10 * time.Second)

	// x, _ := rand.Int(rand.Reader, max)
	// y, _ := rand.Int(rand.Reader, max)
	// z, _ := rand.Int(rand.Reader, max)

	// currentPosition := &fbs.PositionT{
	// 	X: float32(x.Int64()),
	// 	Y: float32(y.Int64()),
	// 	Z: float32(z.Int64()),
	// }

	currentPosition := &fbs.PositionT{
		X: 0,
		Y: 0,
		Z: 0,
	}

end:
	for {
		select {
		case <-after:
			break end
		case <-tick.C:
			// currentPosition = &fbs.PositionT{
			// 	X: currentPosition.X + 1,
			// 	Y: currentPosition.Y + 1,
			// 	Z: currentPosition.Z + 1,
			// }
			{
				builder := flatbuffers.NewBuilder(2000)
				u := fbs.UserT{
					Name: name,
					Pos:  currentPosition,
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
				n, err := io.ReadFull(stream, lengthBuf[:])
				if err != nil {
					return fmt.Errorf("failed to read length: %w\n", err)
				}
				if n != 4 {
					fmt.Printf("failed to read length. not 4")
				}
				length := binary.LittleEndian.Uint32(lengthBuf[:])
				fmt.Printf("[%s]%d\n", name, length)

				buf := make([]byte, length)
				n, err = io.ReadFull(stream, buf)
				if err != nil {
					return fmt.Errorf("[%s]failed to read message: %w", name, err)
				}
				if uint32(n) != length {
					fmt.Printf("[%s]failed to read length. got=%d, want %d\n", name, n, length)
				}
				bs := fbs.GetRootAsBroadcast(buf, 0)
				_ = bs.UnPack()

				// var sb strings.Builder

				// sb.WriteString(fmt.Sprintf("%s: elem=%d, length=%d\n", name, len(bb.Poss), length))
				// for _, u := range bb.Poss {
				// 	if u == nil {
				// 		continue
				// 	}
				// 	if u.Pos == nil {
				// 		sb.WriteString(fmt.Sprintf("\t[%s]no pos\n", u.Name))
				// 	} else {
				// 		sb.WriteString(fmt.Sprintf("\t[%s]{x,y,z}={%f,%f,%f}\n", u.Name, u.Pos.X, u.Pos.Y, u.Pos.Z))
				// 	}
				// }
				// fmt.Println(sb.String() + "end")
			}
		}
	}
	return nil
}
