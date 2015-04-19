package worker

import (
	"fmt"
	"github.com/jacoelho/docker-notifier/notifier"
	"github.com/jacoelho/docker-notifier/notifier/slack"
	"sync"

	dockerapi "github.com/fsouza/go-dockerclient"
)

type Worker struct {
	sync.Mutex
	docker     *dockerapi.Client
	Containers map[string]string
	Alert      notifier.Notifier
}

func New(docker *dockerapi.Client) *Worker {
	return &Worker{
		docker:     docker,
		Containers: make(map[string]string),
		Alert:      slack.New("http://google.com", "cenas", "cenas"),
	}
}

func (w *Worker) Add(containerId string) {
	w.Lock()
	defer w.Unlock()

	container, _ := w.docker.InspectContainer(containerId)
	name := container.Name[1:]

	w.Containers[containerId] = name
	w.Alert.Notify(fmt.Sprintf("%s up", name))
}

func (w *Worker) Remove(containerId string) {
	//inspect docker
	w.Lock()
	defer w.Unlock()

	containerName := w.Containers[containerId]
	delete(w.Containers, containerId)
	w.Alert.Notify(fmt.Sprintf("%s down", containerName))
}
