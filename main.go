package main

import (
	//"fmt"
	dockerapi "github.com/fsouza/go-dockerclient"
	"github.com/jacoelho/docker-notifier/worker"
	"log"
)

func main() {
	docker, err := dockerapi.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		log.Fatal(err)
	}

	events := make(chan *dockerapi.APIEvents)
	if err != nil {
		log.Fatal(err)
	}

	w := worker.New(docker)
	go w.Add("teste")
	go w.Add("teste")

	quit := make(chan struct{})

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
