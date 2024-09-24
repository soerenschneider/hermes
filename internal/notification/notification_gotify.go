package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Gotify struct {
	baseUrl string
	token   string
	client  HttpClient
}

func NewGotify(baseUrl string, token string, client HttpClient) (*Gotify, error) {
	if len(baseUrl) == 0 {
		return nil, errors.New("empty baseUrl provided")
	}

	if len(token) == 0 {
		return nil, errors.New("empty token provided")
	}

	if client == nil {
		client = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	return &Gotify{
		baseUrl: baseUrl,
		client:  client,
		token:   token,
	}, nil
}

func (a *Gotify) Send(ctx context.Context, subject, message string) error {
	request, err := a.getRequest(ctx, message, subject)
	if err != nil {
		return err
	}

	resp, err := a.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (a *Gotify) getRequest(ctx context.Context, subject, text string) (*http.Request, error) {
	baseURL := fmt.Sprintf("%s/message", a.baseUrl)
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	query := url.Query()
	query.Add("token", a.token)
	url.RawQuery = query.Encode()

	data := map[string]any{
		"message":  text,
		"title":    subject,
		"priority": 5,
	}
	marshalled, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), bytes.NewReader(marshalled))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return req, nil
}
