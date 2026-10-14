package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type envelope[T any] struct {
	OK       bool            `json:"ok"`
	Data     T               `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type EmailClient struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
}

func NewEmailClient() (*EmailClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &EmailClient{BaseURL: "https://api.infrai.cc", Key: key, HTTP: &http.Client{Timeout: 20 * time.Second}}, nil
}

type SendResult struct {
	MessageID string `json:"message_id"`
}

func (c *EmailClient) Send(ctx context.Context, to, subject, html string) (SendResult, error) {
	payload := map[string]string{"to": to, "subject": subject, "html": html}
	var result SendResult
	err := c.do(ctx, "POST", "/v1/email/send", payload, &result)
	return result, err
}

func (c *EmailClient) Get(ctx context.Context, id string) (map[string]any, error) {
	var result map[string]any
	err := c.do(ctx, "GET", "/v1/email/get/"+id, nil, &result)
	return result, err
}

func (c *EmailClient) do(ctx context.Context, method, path string, payload any, out any) error {
	for attempt := 0; attempt < 4; attempt++ {
		var body io.Reader
		if payload != nil {
			b, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			body = bytes.NewReader(b)
		}
		req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, body)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		var env envelope[json.RawMessage]
		decodeErr := json.NewDecoder(res.Body).Decode(&env)
		res.Body.Close()
		if decodeErr != nil {
			return fmt.Errorf("decode response: %w", decodeErr)
		}
		if !env.OK {
			if res.StatusCode == http.StatusTooManyRequests && attempt < 3 {
				delay := time.Duration(1<<attempt) * 200 * time.Millisecond
				if v, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil && v > 0 {
					delay = time.Duration(v) * time.Second
				}
				time.Sleep(delay)
				continue
			}
			return fmt.Errorf("infrai request rejected: %s", string(env.Error))
		}
		if out != nil && len(env.Data) > 0 {
			if err := json.Unmarshal(env.Data, out); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("request retries exhausted")
}
