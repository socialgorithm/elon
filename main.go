package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/socialgorithm/elon-server/domain"
)

var upgrader = websocket.Upgrader{}

func main() {
	log.Println("Starting Elon Server")
	http.HandleFunc("/drive", newGameHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer connection.Close()

	carStateUpdates := createCarStateChannel()
	//carControlUpdates := createCarControlStateChannel()
	go writeStateUpdatesToConnection(carStateUpdates, connection)
	//go writeControlUpdatesToEngine(connection, carControlUpdates)
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

func createCarStateChannel() <-chan domain.CarState {
	carStateChannel := make(chan domain.CarState)
	log.Println("Start car state generator")
	go func() {
		for {
			log.Println("Generating new car state")
			carStateChannel <- genRandomCarState()
		}
	}()

	return carStateChannel
}

func genRandomCarState() domain.CarState {
	return domain.CarState{
		Position: domain.Position{X: 1.0, Y: 1.0},
		Velocity: 1,
		Sensors: []domain.Sensor{
			domain.Sensor{Angle: 1, Distance: 1},
		},
	}
}
