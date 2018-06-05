package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/socialgorithm/elon-server/render"
	"github.com/socialgorithm/elon-server/simulator"
)

var upgrader = websocket.Upgrader{}
var simulation *simulator.Simulation

func main() {
	port := "8080"
	simulation = simulator.CreateSimulation(1)

	log.Printf("Starting Elon Server on localhost:%s", port)
	http.HandleFunc("/", connectionHandler)
	go http.ListenAndServe(":"+port, nil)

	render.Render(simulation)
}

func connectionHandler(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer connection.Close()

	for {
		_, data, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error:", err)
			break
		}

		messageStr := string(data[:])
		message := strings.Split(messageStr, " ")

		switch message[0] {
		case "start":
			go simulation.Start()
		}
	}
}
