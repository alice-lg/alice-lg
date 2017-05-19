package birdwatcher

// Http Birdwatcher Client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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
func (self *Client) GetJson(endpoint string) (ClientResponse, error) {
	res, err := http.Get(self.Api + endpoint)
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
