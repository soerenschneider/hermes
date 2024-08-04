package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type Awtrix struct {
	host   string
	client Client
}

func NewAwtrix(host string, client Client) (*Awtrix, error) {
	if client == nil {
		client = http.DefaultClient
	}

	return &Awtrix{
		host:   host,
		client: client,
	}, nil
}

func (a *Awtrix) Send(ctx context.Context, subject, message string) error {
	request, err := a.getRequest(message)
	if err != nil {
		return err
	}

	_, err = a.client.Do(request)
	if err != nil {
		return err
	}

	return nil
}

func (a *Awtrix) getRequest(text string) (*http.Request, error) {
	url := fmt.Sprintf("%s/api/notify", a.host)
	data := map[string]any{
		"text":     text,
		"rainbow":  true,
		"duration": 10,
	}

	marshalled, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshalled))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return req, nil
}
