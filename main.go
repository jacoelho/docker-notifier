package main

import (
	//"fmt"
	dockerapi "github.com/fsouza/go-dockerclient"
	"github.com/jacoelho/docker-notifier/worker"
	"log"
)

func main() {
	docker, err := dockerapi.NewClient("unix:///tmp/docker.sock")
	if err != nil {
		log.Fatal(err)
	}

	events := make(chan *dockerapi.APIEvents)
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan struct{})
	w := worker.New(docker)

	// Process Docker events
	for msg := range events {
		switch msg.Status {
		case "start":
			go w.Add(msg.ID)
		case "die":
			go w.Remove(msg.ID)
		case "stop", "kill":
			go w.Remove(msg.ID)
		}
	}

	close(quit)
}
