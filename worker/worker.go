package worker

import (
	"fmt"
	"github.com/jacoelho/docker-notifier/notifier"
	"log"
	"os"
	"sync"

	dockerapi "github.com/fsouza/go-dockerclient"
)

type Worker struct {
	sync.Mutex
	docker     *dockerapi.Client
	Containers map[string]string
	Alert      notifier.Plugin
}

func New(docker *dockerapi.Client, arguments []string) *Worker {

	if _, ok := notifier.AvailableNotifiers[arguments[0]]; !ok {
		fmt.Printf("invalid plugin\n")
		os.Exit(1)
	}

	alert := notifier.AvailableNotifiers[arguments[0]]().(notifier.Plugin)

	alert.Init(arguments)

	return &Worker{
		docker:     docker,
		Containers: make(map[string]string),
		Alert:      alert,
	}
}

func (w *Worker) RegisterRunning() {
	containers, err := w.docker.ListContainers(dockerapi.ListContainersOptions{})
	if err != nil {
		log.Fatalf("Unable to register running containers: %v", err)
	}
	for _, container := range containers {
		w.Lock()
		c, _ := w.docker.InspectContainer(container.ID)
		name := c.Name[1:]

		w.Containers[container.ID] = name
		w.Unlock()
	}
}

func (w *Worker) Add(containerId string) {
	w.Lock()
	defer w.Unlock()

	container, _ := w.docker.InspectContainer(containerId)
	name := container.Name[1:]

	w.Containers[containerId] = name
	w.Alert.NotifyUp(name)
}

func (w *Worker) Remove(containerId string) {
	w.Lock()
	defer w.Unlock()

	containerName := w.Containers[containerId]
	delete(w.Containers, containerId)
	w.Alert.NotifyDown(containerName)
}
