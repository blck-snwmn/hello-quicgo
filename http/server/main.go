package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/blck-snwmn/hello-quicgo/schema/fbs"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// log.Println("recieve")

		var lengthBuf [4]byte
		_, err := r.Body.Read(lengthBuf[:])
		if err != nil {
			if err != io.EOF {
				fmt.Printf("failed to read length: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		// memo: endian が怪しい？？
		length := binary.LittleEndian.Uint32(lengthBuf[:])
		buf := make([]byte, length)
		_, err = r.Body.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("failed to read message: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		_ = fbs.GetRootAsUser(buf, 0)
		// fmt.Printf("name=%s\n", user.Name())

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello, World!")
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
