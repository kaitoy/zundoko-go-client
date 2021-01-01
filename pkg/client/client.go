package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kaitoy/zundoko-go-client/pkg/model"
	"github.com/kaitoy/zundoko-go-client/pkg/util"
)

// Client represents a Zundoko client.
type Client interface {
	// GetZundokos calls GET Zundokos API and returns the results.
	GetZundokos() ([]model.Zundoko, error)

	// PostZundoko calls POST Zundoko API and returns the result.
	PostZundoko(zundoko *model.Zundoko) error

	// PostKiyoshi calls POST Kiyoshi API and returns the result.
	PostKiyoshi(kiyoshi *model.Kiyoshi) error
}

// NewClient creates a Client instance.
func NewClient(urlBase string) Client {
	return &client{
		urlBase,
		&http.Client{
			Timeout: 10 * time.Second,
		},
		model.NewZundokoDecoder(),
	}
}

// client implements Client interface.
type client struct {
	urlBase        string
	httpClient     util.HTTPClient
	zundokoDecoder model.ZundokoDecoder
}

func (c *client) GetZundokos() ([]model.Zundoko, error) {
	req, _ := http.NewRequest("GET", c.urlBase+"/zundokos", nil)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET Zundoko API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("GET Zundoko API returned an error. status: %s", resp.Status)
		return nil, err
	}

	return c.zundokoDecoder.DecodeList(resp.Body)
}

func (c *client) PostZundoko(zundoko *model.Zundoko) error {
	zundokoJSON, _ := json.Marshal(zundoko)
	req, _ := http.NewRequest(
		"POST",
		c.urlBase+"/zundokos",
		bytes.NewBuffer([]byte(zundokoJSON)),
	)
	req.Header.Add("Content-type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("POST Zundoko API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf("POST Zundoko API returned an error. status: %s", resp.Status)
	}

	return nil
}

func (c *client) PostKiyoshi(kiyoshi *model.Kiyoshi) error {
	kiyoshiJSON, _ := json.Marshal(kiyoshi)
	req, _ := http.NewRequest(
		"POST",
		c.urlBase+"/kiyoshies",
		bytes.NewBuffer([]byte(kiyoshiJSON)),
	)
	req.Header.Add("Content-type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("POST Kiyoshi API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf("POST Kiyoshi API returned an error. status: %s", resp.Status)
	}

	return nil
}
