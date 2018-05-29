package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/socialgorithm/elon-server/domain"
	"github.com/socialgorithm/elon-server/render"
	"github.com/socialgorithm/elon-server/simulator"
)

var upgrader = websocket.Upgrader{}
var simulation simulator.Simulation

func main() {
	log.Println("Starting Elon Server")
	simulation = simulator.StartSimulation()
	fmt.Println("done starting simulation")
	render.Render(simulation)
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

	go writeStateUpdatesToConnection(simulation.CarStateEmitter, connection)
	go writeControlUpdatesToReceiver(connection, simulation.CarControlStateReceiver)
}

func writeStateUpdatesToConnection(carStateChannel <-chan domain.CarState, connection *websocket.Conn) {
	for {
		newCarState := <-carStateChannel
		newCarStateJSON, _ := json.Marshal(newCarState)
		log.Printf("Car state changed: %s", newCarStateJSON)
		connection.WriteMessage(websocket.TextMessage, newCarStateJSON)
	}
}

func writeControlUpdatesToReceiver(connection *websocket.Conn, carControlStateChannel chan<- domain.CarControlState) {
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
