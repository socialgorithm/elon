package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/socialgorithm/elon-server/domain"
	"github.com/socialgorithm/elon-server/simulator"
)

var upgrader = websocket.Upgrader{}

func main() {
	log.Println("Starting Elon Server")
	http.HandleFunc("/simulate", newSimulationHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func newSimulationHandler(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer connection.Close()

	carControlUpdates, carStateUpdates, err := simulator.StartSimulation()
	go writeStateUpdatesToConnection(carStateUpdates, connection)
	go writeControlUpdatesToEngine(connection, carControlUpdates)
}

func writeStateUpdatesToConnection(carStateChannel <-chan domain.CarState, connection *websocket.Conn) {
	for {
		newCarState := <-carStateChannel
		newCarStateJSON, _ := json.Marshal(newCarState)
		log.Printf("Car state changed: %s", newCarStateJSON)
		connection.WriteMessage(websocket.TextMessage, newCarStateJSON)
	}
}

func writeControlUpdatesToEngine(connection *websocket.Conn, carControlStateChannel chan<- domain.CarControlState) {
	defer close(carControlStateChannel)
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error:", err)
			break
		}
		log.Printf("Received client message: %s", message)
	}
}
