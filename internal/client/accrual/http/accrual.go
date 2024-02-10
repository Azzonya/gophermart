package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, orderNumber)

	// Создаем новый HTTP-запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			slog.Error("cannot close response body")
		}
	}()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	} else if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("too many requests")
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("something went wrong")
	}

	body, err := io.ReadAll(resp.Body)
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
