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
	port := "8080"
	log.Printf("Starting Elon Server on localhost:%s", port)
	http.HandleFunc("/", connectionHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
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
			carCount, err := strconv.Atoi(message[1])
			if err != nil {
				log.Println("Error:", err)
				break
			}
			simulation = simulator.CreateSimulation(carCount)
			break
		case "start":
			go simulation.Start()
			render.Render(simulation)
		}
	}
}
