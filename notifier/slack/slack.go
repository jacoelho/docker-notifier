package slack

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	. "github.com/jacoelho/docker-notifier/notifier"
)

type SlackIncomingHook struct {
	Channel     string             `json:"channel"`
	Username    string             `json:"username"`
	IconUrl     string             `json:"icon_url"`
	Attachments []SlackAttachments `json:"attachments"`
}

type SlackAttachments struct {
	Fallback string `json:"fallback"`
	Color    string `json:"color"`
	Title    string `json:"title"`
	Text     string `json:"text"`
}

type SlackNotifier struct {
	Username string
	Url      string
	Channel  string
}

func (s *SlackNotifier) Init(parameters []string) {
	flags := flag.NewFlagSet("slack", flag.ExitOnError)
	flags.StringVar(&s.Username, "username", "docker-notifier", "docker user")
	flags.StringVar(&s.Url, "url", "required", "slack incoming hook")
	flags.StringVar(&s.Channel, "channel", "required", "slack channel")
	flags.Parse(parameters[1:])

	if s.Url == "required" {
		fmt.Println("error: enter a slack incoming hook")
		os.Exit(1)
	}

	if s.Channel == "required" {
		fmt.Println("error: enter a channel name")
		os.Exit(1)
	}

	if !(strings.HasPrefix(s.Channel, "#") != strings.HasPrefix(s.Channel, "@")) {
		fmt.Println("invalid channel name")
		os.Exit(1)
	}
}

func (s *SlackNotifier) NotifyUp(containerName string) {
	s.Notify(fmt.Sprintf("container %s is up", containerName), "good")
}

func (s *SlackNotifier) NotifyDown(containerName string) {
	s.Notify(fmt.Sprintf("container %s is down", containerName), "danger")
}

func (s *SlackNotifier) Notify(text string, color string) {
	body := &SlackIncomingHook{
		Channel:  s.Channel,
		Username: s.Username,
		IconUrl:  "https://raw.githubusercontent.com/jacoelho/docker-notifier/master/docker.png",
		Attachments: []SlackAttachments{
			SlackAttachments{
				Fallback: text,
				Color:    color,
				Title:    s.Username,
				Text:     text,
			},
		},
	}

	postJson, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", s.Url, bytes.NewReader(postJson))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("error post: %v\n", resp.StatusCode)
	}
}

func init() {
	RegisterNotifier("slack", func() interface{} {
		return new(SlackNotifier)
	})
}
