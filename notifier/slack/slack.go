package slack

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type SlackIncomingHook struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Text     string `json:"text"`
	IconUrl  string `json:"icon_url"`
}

type Notifier struct {
	Url string
}

func New(url string) *Notifier {
	return &Notifier{
		Url: url,
	}
}

func (s *Notifier) Notify(text string) {
	body := &SlackIncomingHook{
		Channel:  "#integration-test",
		Username: "docker-notifier",
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
