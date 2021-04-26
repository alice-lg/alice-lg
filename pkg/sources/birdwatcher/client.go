package birdwatcher

// Http Birdwatcher Client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type ClientResponse map[string]interface{}

type Client struct {
	Api string
}

func NewClient(api string) *Client {
	client := &Client{
		Api: api,
	}
	return client
}

// Make API request, parse response and return map or error
func (self *Client) Get(client *http.Client, url string) (ClientResponse, error) {
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

// Make API request, parse response and return map or error
func (self *Client) GetJson(endpoint string) (ClientResponse, error) {
	client := &http.Client{}

	return self.Get(client, self.Api+endpoint)
}

// Make API request, parse response and return map or error
func (self *Client) GetJsonTimeout(timeout time.Duration, endpoint string) (ClientResponse, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	return self.Get(client, self.Api+endpoint)
}
