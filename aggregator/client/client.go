package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
)

type Client struct {
	EndPoint string
}

func NewClient(endPoint string) *Client {
	return &Client{
		EndPoint: endPoint,
	}
}

func (c *Client) AggregateInvoice(dist *types.Distance) error {
	b, err := json.Marshal(dist)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.EndPoint, bytes.NewReader(b))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non http status code 200, got %d", resp.StatusCode)
	}
	return nil
}
