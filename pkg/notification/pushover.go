package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	PushoverMessageURL = "https://api.pushover.net/1/messages.json"
)

type Pushover struct {
	userKey string
	apiKey  string
}

type PushoverMessage struct {
	Token   string `json:"token"`
	User    string `json:"user"`
	Message string `json:"message"`
	Title   string `json:"title"`
}

func NewPushover(userKey string, apiKey string) Notificator {
	return &Pushover{
		userKey: userKey,
		apiKey:  apiKey,
	}
}

func (p Pushover) Name() string {
	return "Pushover"
}

func (p Pushover) Notify(title string, message string) error {
	msg := PushoverMessage{
		Token:   p.apiKey,
		User:    p.userKey,
		Message: message,
		Title:   title,
	}

	out, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	res, err := http.Post(PushoverMessageURL, "application/json", bytes.NewReader(out))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("failed to send notification")
	}

	return nil
}
