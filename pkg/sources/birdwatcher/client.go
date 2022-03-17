package birdwatcher

// Http Birdwatcher Client

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// ClientResponse is a key value mapping
type ClientResponse map[string]interface{}

// A Client uses the http client to talk
// to the birdwatcher API.
type Client struct {
	api string
}

// NewClient creates a new client instance
func NewClient(api string) *Client {
	client := &Client{
		api: api,
	}
	return client
}

// GetEndpoint makes an API request and returns the
// response. The response body will be parsed further
// downstream.
func (c *Client) GetEndpoint(
	ctx context.Context,
	endpoint string,
) (*http.Response, error) {
	client := &http.Client{}
	return client.Get(c.api + endpoint)
}

// Get makes an API request.
// Parse response and return map or error.
func (c *Client) Get(
	client *http.Client,
	url string,
) (ClientResponse, error) {
	res, err := client.Get(url)
	if err != nil {
		return ClientResponse{}, err
	}

	// Read body
	defer res.Body.Close()
	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ClientResponse{}, err
	}

	// Decode json payload
	result := make(ClientResponse)
	err = json.Unmarshal(payload, &result)
	if err != nil {
		return ClientResponse{}, err
	}

	return result, nil
}

// GetJSON makes an API request.
// Parse JSON response and return map or error.
func (c *Client) GetJSON(
	endpoint string,
) (ClientResponse, error) {
	client := &http.Client{}
	return c.Get(client, c.api+endpoint)
}

// GetJSONTimeout make an API request, parses the
// JSON response and returns the result or an error.
//
// This will all go away one we use the new context.
func (c *Client) GetJSONTimeout(
	timeout time.Duration,
	endpoint string,
) (ClientResponse, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	return c.Get(client, c.api+endpoint)
}
