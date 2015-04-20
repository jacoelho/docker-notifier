package slack

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	. "notifier"
)

type SlackIncomingHook struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Text     string `json:"text"`
	IconUrl  string `json:"icon_url"`
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
}

func (s *SlackNotifier) Notify(text string) {
	body := &SlackIncomingHook{
		Channel:  "#integration-test",
		Username: s.Username,
		Text:     text,
		IconUrl:  "https://raw.githubusercontent.com/jacoelho/docker-notifier/master/docker.png",
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
