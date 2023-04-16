package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/blck-snwmn/hello-quicgo/schema/fbs"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

func main() {
	s := webtransport.Server{
		H3: http3.Server{Addr: ":4433"},
	}

	chats := &chat{
		rooms: map[string]*room{},
	}

	// Create a new HTTP endpoint /webtransport.
	http.HandleFunc("/webtransport", func(w http.ResponseWriter, r *http.Request) {
		log.Println("recieve")
		defer log.Println("done")

		// リクエストからルーム名を取得
		roomName := r.URL.Query().Get("room")

		conn, err := s.Upgrade(w, r)
		if err != nil {
			log.Printf("upgrading failed: %s", err)
			w.WriteHeader(500)
			return
		}
		if err := chats.handleSession(conn, roomName); err != nil {
			log.Printf("handling session failed: %s", err)
			w.WriteHeader(500)
			return
		}
	})

	err := s.ListenAndServeTLS("./testdata/server.crt", "./testdata/server.key")
	if err != nil {
		log.Fatal(err)
	}
}

type user struct {
	name     string
	position position
}

type position struct {
	x float32
	y float32
	z float32
}

type room struct {
	positions map[string]position
	mux       sync.Mutex
}

func (rm *room) updatePosition(name string, pos position) {
	rm.mux.Lock()
	rm.positions[name] = pos
	rm.mux.Unlock()
}

// getPositions returns a slice of positions.
func (rm *room) getUserPositions() []user {
	rm.mux.Lock()
	defer rm.mux.Unlock()
	users := make([]user, 0, len(rm.positions))
	for name, pos := range rm.positions {
		users = append(users, user{name: name, position: pos})
	}
	return users

}

type chat struct {
	rooms map[string]*room
	mux   sync.Mutex
}

func (r *chat) getRoom(name string) *room {
	r.mux.Lock()
	rm, ok := r.rooms[name]
	if !ok {
		rm = &room{positions: map[string]position{}}
		r.rooms[name] = rm
	}
	r.mux.Unlock()
	return rm
}

func (ch *chat) handleSession(conn *webtransport.Session, roomName string) error {
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return fmt.Errorf("opening stream failed: %w", err)
	}
	for {
		var lengthBuf [4]byte
		_, err = stream.Read(lengthBuf[:])
		if err != nil {
			return fmt.Errorf("failed to read length: %w", err)
		}
		length := binary.LittleEndian.Uint32(lengthBuf[:])

		buf := make([]byte, length)
		_, err = stream.Read(buf)
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		user := fbs.GetRootAsUser(buf, 0)

		var pos fbs.Position
		user.Pos(&pos)

		room := ch.getRoom(roomName)
		room.updatePosition(string(user.Name()), position{
			x: pos.X(),
			y: pos.Y(),
			z: pos.Z(),
		})

		positions := room.getUserPositions()

		b := fbs.BroadcastT{
			Poss: make([]*fbs.UserPositionT, 0, len(positions)),
		}
		for _, u := range positions {
			up := fbs.UserPositionT{
				Name: u.name,
				Pos: &fbs.PositionT{
					X: u.position.x,
					Y: u.position.y,
					Z: u.position.z,
				},
			}
			b.Poss = append(b.Poss, &up)
		}
		builder := flatbuffers.NewBuilder(200)
		builder.FinishSizePrefixed(b.Pack(builder))
		buf = builder.FinishedBytes()

		_, err := stream.Write(buf)
		if err != nil {
			fmt.Printf("failed to wirte message: %v", err)
		}
	}
}
