package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/socialgorithm/elon-server/clientrenderer"
	"github.com/socialgorithm/elon-server/domain"
)

var inputs = make(chan domain.CarControlState)

func update(connection *websocket.Conn) {
	for {
		select {
		case carControlState := <-inputs:
			message := fmt.Sprintf("input %f %f", carControlState.Steering, carControlState.Throttle)
			connection.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080"}
	log.Printf("connecting to %s", u.String())
	test := true

	connection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer connection.Close()

	log.Println("Write start message")
	connection.WriteMessage(websocket.TextMessage, []byte("start"))

	go update(connection)

	if test {
		clientrenderer.Manual(inputs)
	} else {
		// clientrenderer.ExternalProcess()
	}

	log.Println("Closing connection")
	err = connection.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)

	if err != nil {
		log.Println("Close Error:", err)
		return
	}
}
