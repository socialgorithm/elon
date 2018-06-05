package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/socialgorithm/elon-server/render"
	"github.com/socialgorithm/elon-server/simulator"
)

var upgrader = websocket.Upgrader{}
var simulation simulator.Simulation

func main() {
	log.Println("Starting Elon Server")
	simulation = simulator.CreateSimulation(5)
	go simulation.Start()

	http.HandleFunc("/", connectionHandler)
	go http.ListenAndServe(":8080", nil)

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
		case "init":
			_, err := strconv.Atoi(message[1])
			if err != nil {
				log.Println("Error:", err)
				break
			}
			break
		case "start":
			go simulation.Start()
		}
	}
}
