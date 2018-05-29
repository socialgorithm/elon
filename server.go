package main

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/socialgorithm/elon-server/render"
	"github.com/socialgorithm/elon-server/simulator"
)

var upgrader = websocket.Upgrader{}
var simulation simulator.Simulation

func main() {
	log.Println("Starting Elon Server")
	simulation = simulator.PrepareSimulation()
	go func() {
		start()
	}()
	render.Render(simulation)
}

// Start the server here
func start() {
	simulator.StartSimulation(simulation)
	// http.HandleFunc("/simulate", newSimulationHandler)
	// log.Fatal(http.ListenAndServe(":8080", nil))
}

// physics.runTick(track Track, car Car) car Car {

// }

// func newSimulationHandler(w http.ResponseWriter, r *http.Request) {
// 	connection, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer connection.Close()

// 	simulator.StartSimulation(simulation)

// 	go writeStateUpdatesToConnection(simulation.CarStateChannel, connection)
// 	go writeControlUpdatesToReceiver(connection, simulation.CarControlStateReceiver)
// }

// func writeStateUpdatesToConnection(carStateChannel <-chan domain.CarState, connection *websocket.Conn) {
// 	for {
// 		newCarState := <-carStateChannel
// 		newCarStateJSON, _ := json.Marshal(newCarState)
// 		log.Printf("Car state changed: %s", newCarStateJSON)
// 		connection.WriteMessage(websocket.TextMessage, newCarStateJSON)
// 	}
// }

// func writeControlUpdatesToReceiver(connection *websocket.Conn, carControlStateChannel chan<- domain.CarControlState) {
// 	defer close(carControlStateChannel)
// 	for {
// 		_, message, err := connection.ReadMessage()
// 		if err != nil {
// 			log.Println("Error:", err)
// 			break
// 		}
// 		log.Printf("Received client message: %s", message)
// 	}
// }
