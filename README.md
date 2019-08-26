# Elon Self Driving Car Simulator - Server

This is the server for Socialgorithm's self-driving car simulator, featuring socket listeners that allow clients to connect, receive sensor data and send control data.

## Running

Make sure you have installed at least go 1.11, as that's the version that introduced support for go modules.
Dependencies for the server will be automatically installed when you start it for the first time.

To start the server:

```
$ go run main.go
```

To run the test client run:

```
$ go run client/main.go
```

By default the server runs in "test" mode, which spawns a simulation with 1 car and starts running the simulation clock.
