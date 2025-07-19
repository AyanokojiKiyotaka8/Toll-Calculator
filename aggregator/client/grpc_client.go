package client

import (
	"context"
	"fmt"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	EndPoint string
	types.AggregatorClient
	Conn *grpc.ClientConn
}

func NewGRPCClient(endPoint string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(endPoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %v", err)
	}
	return &GRPCClient{
		EndPoint:         endPoint,
		AggregatorClient: types.NewAggregatorClient(conn),
		Conn:             conn,
	}, nil
}

func (c *GRPCClient) Aggregate(ctx context.Context, aggReq *types.AggregatorReq) error {
	_, err := c.AggregatorClient.Aggregate(ctx, aggReq)
	return err
}

func (c *GRPCClient) Close() error {
	return c.Conn.Close()
}
