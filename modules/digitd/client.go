package digitd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	Endpoint string
	Timeout  time.Duration
}

func NewClient(ip string) *Client {
	return &Client{
		Endpoint: fmt.Sprintf("http://%s:8080/api/learn", ip),
		Timeout:  5 * time.Second,
	}
}

// Speak relays a message to the DigitD voice node
func (c *Client) Speak(topic, content string) error {
	payload := map[string]string{
		"topic":   topic,
		"content": content,
	}
	jsonData, _ := json.Marshal(payload)

	httpClient := &http.Client{Timeout: c.Timeout}
	resp, err := httpClient.Post(c.Endpoint, "application/json", bytes.NewBuffer(jsonData))
	
	if err != nil {
		return fmt.Errorf("digitd unreachable: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("digitd rejected message: status %d", resp.StatusCode)
	}
	return nil
}
