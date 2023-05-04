package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/blck-snwmn/hello-quicgo/schema/fbs"
	flatbuffers "github.com/google/flatbuffers/go"
)

func main() {
	i := 0
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

	var sg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		sg.Add(1)
		go func() {
			defer sg.Done()
			resp, err := http.Post("http://localhost:8080/", "application/octet-stream", bytes.NewReader(buf))
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(body))
		}()
	}
	sg.Wait()
}
