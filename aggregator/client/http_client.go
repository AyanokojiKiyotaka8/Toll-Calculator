package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
)

type HTTPClient struct {
	EndPoint string
}

func NewHTTPClient(endPoint string) *HTTPClient {
	return &HTTPClient{
		EndPoint: endPoint,
	}
}

func (c *HTTPClient) Aggregate(ctx context.Context, aggReq *types.AggregatorReq) error {
	dist := &types.Distance{
		OBUID: int(aggReq.GetObuId()),
		Value: aggReq.GetValue(),
		Unix:  aggReq.GetUnix(),
	}

	b, err := json.Marshal(dist)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.EndPoint, bytes.NewReader(b))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non http status code 200, got %d", resp.StatusCode)
	}
	return nil
}
