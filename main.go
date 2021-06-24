package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/socialgorithm/elon-server/domain"

	"github.com/gorilla/websocket"
	"github.com/socialgorithm/elon-server/render"
	"github.com/socialgorithm/elon-server/simulator"
)

var upgrader = websocket.Upgrader{}
var simulation *simulator.Simulation
var isTest = false

func main() {
	var port = flag.String("port", "8080", "the port number to run on")
	var test = flag.Bool("test", true, "whether the server should run in test mode")
	isTest = *test
	simulation = simulator.CreateSimulation(1)

	if isTest {
		go simulation.Start(isTest)
	}

	log.Printf("Starting Elon Server on localhost:%s", *port)
	http.HandleFunc("/", connectionHandler)
	go http.ListenAndServe(":"+*port, nil)

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
			go simulation.Start(isTest)
			break
		case "input":
			steering, _ := strconv.ParseFloat(message[1], 64)
			throttle, _ := strconv.ParseFloat(message[2], 64)

			carControlState := domain.CarControlState{
				Steering: steering,
				Throttle: throttle,
			}
			go simulation.Input(0, carControlState)
			break
		case "control":
			signal, err := strconv.Atoi(message[1])
			if err != nil {
				fmt.Printf("Received invalid signal %s\n", message[1])
				break
			}
			if signal == simulator.SimulationRestart {
				simulation.Restart()
			}
		}
	}
}
