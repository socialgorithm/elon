# Elon Self Driving Car Simulator - Server

This is the server for Socialgorithm's self-driving car simulator, featuring socket listeners that allow clients to connect, receive sensor data and send control data.

## Extra dependencies

Some dependencies don't like `dep` too much, so we need to install them manually (for now):

```console
go get github.com/go-gl/glfw/v3.2/glfw
go get github.com/go-gl/gl/v3.3-core/gl
go get github.com/go-gl/mathgl/mgl32
```

### GEOS

You'll need the C library GEOS (used for some geometry calculations)

After installing it (see below), get `gogeos`:

```
$ go get github.com/paulsmith/gogeos/geos
```

#### Ubuntu

```
$ apt-get install libgeos-dev
```
#### OS X - homebrew

```
$ brew install geos
```

## Running

To start the server:

`go run main.go`

By default the server runs in "test" mode, which spawns a simulation with 1 car and starts running the simulation clock.