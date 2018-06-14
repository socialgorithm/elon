package simulator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/socialgorithm/elon-server/domain"
)

const trackDir = "tracks/"

// ReadTrack reads a random track from the tracks folder
func ReadTrack() domain.Track {
	files, err := ioutil.ReadDir(trackDir)
	if err != nil || len(files) < 1 {
		fmt.Println("No track files available")
		os.Exit(-1)
	}

	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	trackFile := files[rand.Intn(len(files))]

	// read track
	trackJSON, err := ioutil.ReadFile(trackDir + trackFile.Name())
	check(err)

	var track domain.Track
	err = json.Unmarshal(trackJSON, &track)
	check(err)

	return track
}

func check(e error) {
	if e != nil {
		fmt.Println("Error reading track file")
		panic(e)
	}
}
