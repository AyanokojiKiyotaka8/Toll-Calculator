package client

import (
	"context"
	"log"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	EndPoint string
	types.AggregatorClient
}

func NewGRPCClient(endPoint string) *GRPCClient {
	conn, err := grpc.NewClient(endPoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return &GRPCClient{
		EndPoint:         endPoint,
		AggregatorClient: types.NewAggregatorClient(conn),
	}
}

func (c *GRPCClient) Aggregate(ctx context.Context, aggReq *types.AggregatorReq) error {
	_, err := c.AggregatorClient.Aggregate(ctx, aggReq)
	return err
}
