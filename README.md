# Elon Self Driving Car Simulator - Server

This is the server for Socialgorithm's self-driving car simulator, featuring socket listeners that allow clients to connect, receive sensor data and send control data.

## Extra dependencies

Some dependencies don't like `dep` too much, so we need to install them manually (for now):

```console
go get github.com/go-gl/glfw/v3.2/glfw
go get github.com/go-gl/gl/v3.3-core/gl
go get github.com/go-gl/mathgl/mgl32
```

## Running

To start the server:

```
$ go run main.go
```

To run the test client run:

```
$ go run client/main.go
```

By default the server runs in "test" mode, which spawns a simulation with 1 car and starts running the simulation clock.