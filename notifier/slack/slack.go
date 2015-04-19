package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SlackIncomingHook struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Text     string `json:"text"`
	Icon     string `json:"icon_emoji"`
}

type Notifier struct {
	Url      string
	Username string
	Icon     string
}

func New(url string, user string, icon string) *Notifier {
	return &Notifier{
		Url:      url,
		Username: user,
		Icon:     icon,
	}
}

func (s *Notifier) Notify(text string) {
	body := &SlackIncomingHook{
		Channel:  "cenas",
		Username: "cenas",
		Text:     text,
		Icon:     "cenas",
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
		fmt.Printf("error post: %v\n", resp.StatusCode)
	}

	fmt.Printf("slack -> %s\n", text)
}
