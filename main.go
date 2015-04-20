package main

import (
	"fmt"
	dockerapi "github.com/fsouza/go-dockerclient"
	"log"
	//	"notifier"
	_ "notifier/slack"
	"os"
	"worker"
)

func getopt(name, fallback string) string {
	if env := os.Getenv(name); env != "" {
		return env
	}
	return fallback
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("wrong number of arguments\n")
		os.Exit(1)
	}

	docker, err := dockerapi.NewClient(getopt("DOCKER_HOST", "unix:///var/run/docker.sock"))
	if err != nil {
		log.Fatal(err)
	}

	events := make(chan *dockerapi.APIEvents)
	if err != nil {
		log.Fatal(err)
	}

	docker.AddEventListener(events)

	w := worker.New(docker, os.Args[1:])

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
