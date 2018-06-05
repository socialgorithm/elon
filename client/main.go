package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080"}
	log.Printf("connecting to %s", u.String())

	connection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer connection.Close()

	connection.WriteMessage(websocket.TextMessage, []byte("init 5"))

	connection.WriteMessage(websocket.TextMessage, []byte("start"))
}
