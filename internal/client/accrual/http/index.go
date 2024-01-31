package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	client  http.Client
	baseURL string
}

func New(uri string) *Client {
	return &Client{
		client:  http.Client{},
		baseURL: uri,
	}
}

func (c *Client) Send(orderNumber string) (*RequestResult, error) {
	url := c.baseURL + fmt.Sprintf("api/orders/%s", orderNumber)

	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent {
		return nil, nil
	} else if response.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("too many requests")
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("something went wrong")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	result := RequestResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
